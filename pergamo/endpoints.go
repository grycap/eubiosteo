package pergamo

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"pergamo/pergamo/structs"
	"pergamo/pergamo/types"
	"strings"
	"time"
)

func addFileByteToZip(name string, content []byte, w *zip.Writer) error {
	header := &zip.FileHeader{
		Name:         name,
		Method:       zip.Store,
		ModifiedTime: uint16(time.Now().UnixNano()),
		ModifiedDate: uint16(time.Now().UnixNano()),
	}

	fw, err := w.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("Cannot create output header %s", err.Error())
	}

	if _, err = io.Copy(fw, bytes.NewReader(content)); err != nil {
		return fmt.Errorf("Cannot copy %s", err.Error())
	}

	return nil
}

func createFileFromJSON(name string, output structs.JSONGeneric, w *zip.Writer) error {
	outputBytes, err := output.ToByte()
	if err != nil {
		return fmt.Errorf("Cannot get byte representation of output")
	}

	return addFileByteToZip(name+".json", outputBytes, w)
}

func addFileToZip(path, filetype, filename string, w *zip.Writer) error {
	data, err := ioutil.ReadFile(filepath.Join(path, filename))
	if err != nil {
		return err
	}

	extension, err := types.GetFileExtension(filetype)
	if err != nil {
		return fmt.Errorf("Failed to get extension from typename %s with error %s", filetype, err.Error())
	}

	return addFileByteToZip(filename+"."+extension, data, w)
}

func (s *Server) Download(req *structs.AllocDownloadRequest, resp *structs.AllocDownloadResponse) error {

	alloc, err := s.stateStore.AllocByID(req.Allocid)
	if err != nil {
		return err
	}

	job, err := s.stateStore.JobByID(alloc.JobID)
	if err != nil {
		return err
	}

	// Add the output file
	err = createFileFromJSON("output", alloc.Output, req.Writer)
	if err != nil {
		return fmt.Errorf("Failed to create output: %s", err.Error())
	}

	// Add the input file
	err = createFileFromJSON("input", alloc.Input, req.Writer)
	if err != nil {
		return fmt.Errorf("Failed to create output: %s", err.Error())
	}

	// load output files
	for outputName, outputType := range job.Output {
		if types.IsFileType(outputType) {
			path, ok := alloc.Output[outputName]
			if !ok {
				return fmt.Errorf("Not found outputname %s on alloc", outputName)
			}

			filename, ok := path.(string)
			if !ok {
				return fmt.Errorf("Cannot cast outputname %s to string", outputName)
			}

			// copy onedata+filepath to writter
			err = addFileToZip(s.config.OnedataPath, outputType, filename, req.Writer)
			if err != nil {
				return fmt.Errorf("Cannot save file to zip with name %s error %s", filename, err.Error())
			}
		}
	}

	return nil
}

func (s *Server) PostAlloc(req *structs.AllocPostRequest, resp *structs.AllocPostResponse) error {
	id, err := s.CreateAlloc(&req.Alloc)
	if err != nil {
		alloc := structs.Alloc{
			ID:     id,
			Status: structs.FAILED,
			Error:  err.Error(),
		}

		err = s.stateStore.UpdateAlloc(&alloc)
		if err != nil {
			fmt.Println("Cannot update alloc with error")
		}

		al, err := s.stateStore.AllocByID(id)
		if err != nil {
			panic(err)
		}

		fmt.Println(al)
	}

	resp.ID = id
	return err
}

func getFilename(content string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("Is empty the header")
	}

	nameOne := strings.Split(content, "filename=\"")
	if len(nameOne) == 2 {
		return strings.Replace(nameOne[1], "\"", "", -1), nil
	}

	return "", fmt.Errorf("length should be two but is %d", len(nameOne))
}

func writeFileToServer(path string, file *multipart.FileHeader) error {
	// copy file into folder
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	content, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, content, 0644)
	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

const tempZipFile = "/tmp/temp.zip"

func (s *Server) PreJobCreateUploadZip(file *multipart.FileHeader) (string, error) {

	if file == nil {
		return "", nil
	}

	filetype := file.Header.Get("Content-Type")
	uid := structs.GenerateUUID()

	filename, err := getFilename(file.Header.Get("Content-Disposition"))
	if err != nil {
		return "", fmt.Errorf("Cannot get filename %s", err.Error())
	}

	fmt.Println("file name")
	fmt.Println(filename)

	if filetype == "application/zip" {
		// unzip folder on storage with another name

		// guardar en tmp
		err = writeFileToServer(tempZipFile, file)
		if err != nil {
			return "", err
		}

		err = Unzip(tempZipFile, filepath.Join(s.config.OnedataPath, uid))
		if err != nil {
			return "", err
		}

		err := os.Remove(tempZipFile)
		if err != nil {
			return "", err
		}

	} else {
		// create folder and add value
		err := os.MkdirAll(filepath.Join(s.config.OnedataPath, uid), 0777)
		if err != nil {
			return "", err
		}

		err = writeFileToServer(filepath.Join(s.config.OnedataPath, uid, filename), file)
		if err != nil {
			return "", err
		}

	}

	return uid, nil
}

func (s *Server) PostJob(req *structs.JobPostRequest, resp *structs.JobPostResponse) error {

	zipPath, err := s.PreJobCreateUploadZip(req.File)
	if err != nil {
		return err
	}

	req.Job.Zip = zipPath

	err = s.CreateJob(&req.Job)
	resp.ID = req.Job.ID
	return err
}

func (s *Server) PostWorkflow(req *structs.WorkflowPostRequest, resp *structs.WorkflowPostResponse) error {

	err := s.CreateWorkflow(req.Workflow)
	resp.ID = req.Workflow.ID

	return err
}

// List endpoints

func (s *Server) PostImage(req *structs.ImagePostRequest, resp *structs.ImagePostResponse) error {

	formatDetected := types.GetImageFormat(&req.Content)
	fmt.Printf("PostImage: Image %s format detected '%s'\n", req.Name, formatDetected)

	if req.Format != "other" && req.Format != formatDetected {
		fmt.Printf("PostImage: Format mismatch between provided %s and detected %s\n", req.Format, formatDetected)
		return fmt.Errorf("Format is %s but found %s", formatDetected, req.Format)
	}

	image := structs.Image{
		Name:   req.Name,
		Format: formatDetected,
	}

	// save image on db
	err := s.stateStore.SaveImage(&image)
	if err != nil {
		return err
	}

	// save image on folder
	err = s.stateStore.SaveImageToStorage(&image, req.Content)
	resp.ID = image.ID

	return err
}

func (s *Server) CheckDownloadImage(id string) (string, string, error) {
	image, err := s.stateStore.GetImage(id)
	if err != nil {
		return "", "", err
	}

	path := filepath.Join(s.config.OnedataPath, id)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", "", fmt.Errorf("Path on storage does not exists for %s on path %s", id, path)
	}

	return path, image.Name, nil
}

func (s *Server) ImageList(req *structs.ImageListRequest, resp *structs.ImageListResponse) error {

	images, hasNext, err := s.stateStore.GetImages(&req.QueryOffset)
	if err != nil {
		return err
	}

	resp.Images = *images
	resp.HasNext = hasNext

	return nil
}

func (s *Server) JobList(req *structs.JobListRequest, resp *structs.JobListResponse) error {

	jobs, hasNext, err := s.stateStore.GetJobs(&req.QueryOffset)
	if err != nil {
		return err
	}

	resp.Jobs = *jobs
	resp.HasNext = hasNext

	return nil
}

func (s *Server) AllocList(req *structs.AllocListRequest, resp *structs.AllocListResponse) error {

	allocs, hasNext, err := s.stateStore.GetAllocs(&req.QueryOffset)
	if err != nil {
		return err
	}

	resp.Allocs = *allocs
	resp.HasNext = hasNext

	return nil
}

func (s *Server) WorkflowList(req *structs.WorkflowListRequest, resp *structs.WorkflowListResponse) error {

	workflows, hasNext, err := s.stateStore.GetWorkflows(&req.QueryOffset)
	if err != nil {
		return err
	}

	resp.Workflows = *workflows
	resp.HasNext = hasNext

	return nil
}

func (s *Server) WorkflowAllocList(req *structs.WorkflowAllocListRequest, resp *structs.WorkflowAllocListResponse) error {

	workflowallocs, hasNext, err := s.stateStore.GetWorkflowAllocs(&req.QueryOffset)
	if err != nil {
		return err
	}

	resp.WorkflowAllocs = *workflowallocs
	resp.HasNext = hasNext

	return nil
}

func (s *Server) WorkflowByID(id string) (*structs.Workflow, error) {
	return s.stateStore.WorkflowByID(id)
}

func (s *Server) ApplyWorkflowAlloc(workflowalloc *structs.WorkflowAlloc) error {

	workflowalloc.Status = structs.RUNNING
	err := s.stateStore.CreateWorkflowAlloc(workflowalloc)
	if err != nil {
		return err
	}

	workflow, err := s.stateStore.WorkflowByID(workflowalloc.WorkflowID)
	if err != nil {
		return err
	}

	jobs, err := s.stateStore.WorkflowJobs(workflow)
	if err != nil {
		return err
	}

	initialStep := workflow.Steps[0]
	alloc, err := CreateAllocFromStep(workflowalloc, &initialStep, jobs, &workflow.Variables, &workflowalloc.Input, nil)
	if err != nil {
		return err
	}

	alloc.Status = structs.RUNNING

	fmt.Println("-- workflow alloc --")
	fmt.Println(alloc)

	// &{74dc7d13-bb96-b091-b36f-180cbcecf976 work0 0 map[entry:4] map[] 0 0 0}

	_, err = s.CreateAlloc(alloc)
	return err
}

func getEntryValues(workflow *structs.Workflow, jobs *[]structs.Job) (structs.JSON, error) {

	mapJobs := map[string]structs.Job{}
	for _, job := range *jobs {
		mapJobs[job.ID] = job
	}

	entries := map[string]string{}

	for _, step := range workflow.Steps {

		job, ok := mapJobs[step.Job]
		if !ok {
			continue // no deberia pasar
		}

		for inputName, inputValue := range step.Param {

			inputValueStr, ok := inputValue.(string)
			if !ok {
				continue
			}

			if !strings.HasPrefix(inputValueStr, "@input_") {
				continue
			}

			inputType, ok := job.Input[inputName]
			if !ok {
				return nil, fmt.Errorf("Buscar input de nombre %s pero no estaba en job %s", inputName, step.Job)
			}

			inputVariableName := strings.Replace(inputValueStr, "@input_", "", -1)
			entries[inputVariableName] = inputType
		}
	}

	return entries, nil
}

func (s *Server) CreateWorkflow(workflow *structs.Workflow) error {

	jobNames := []string{}
	for _, i := range workflow.Steps {
		jobNames = append(jobNames, i.Job)
	}

	fmt.Println(jobNames)

	jobs, err := s.stateStore.JobsByIDs(jobNames)
	if err != nil {
		return err
	}

	if len(workflow.Steps) == 0 {
		return fmt.Errorf("There are no steps on the workflow")
	}

	fmt.Println("-- jobs al validar --")
	fmt.Println(jobs)

	err = ValidateWorkflow(workflow, jobs)
	if err != nil {
		return err
	}

	entryValues, err := getEntryValues(workflow, jobs)
	if err != nil {
		return err
	}

	workflow.Entry = entryValues
	return s.stateStore.CreateWorkflow(workflow)
}

func (s *Server) CreateJob(job *structs.Job) error {
	return s.stateStore.CreateJob(job)
}

func (s *Server) JobByID(id string) (*structs.Job, error) {
	return s.stateStore.JobByID(id)
}

func (s *Server) AllocByID(req *structs.AllocSingleRequest, resp *structs.AllocSingleResponse) error {
	alloc, err := s.stateStore.AllocByID(req.ID)
	resp.Alloc = alloc

	return err
}

func completeAlloc(job *structs.Job, alloc *structs.Alloc) {

	// completar en alloc los atributos que sean de tipo file (output)
	for outputName, outputType := range job.Output {
		if types.IsFileType(outputType) {
			alloc.Output[outputName] = structs.GenerateUUID()
		}
	}

}

func (s *Server) CreateAlloc(alloc *structs.Alloc) (string, error) {

	fmt.Println("-- alloc --")
	fmt.Println(alloc)

	job, err := s.stateStore.JobByID(alloc.JobID)
	if err != nil {
		return "", fmt.Errorf("Alloc. No hay job con id %s", alloc.JobID)
	}

	// No pasamos por el pending...
	alloc.Status = structs.RUNNING

	fmt.Println("-- one data --")
	fmt.Println(s.config.OnedataPath)

	validateErr := ValidateInput(s.config.OnedataPath, job, alloc)
	if validateErr != nil {
		alloc.Error = validateErr.Error()
		alloc.Status = structs.VALIDATIONINPUTFAILURE
	} else {

		// COMPLETAR EL ALLOC
		completeAlloc(job, alloc)
	}

	fmt.Println("-- hola 2 --")
	// createalloc inicializa el alloc con algunas variables importantes (i.e nombres de path)
	err = s.stateStore.CreateAlloc(alloc)
	if err != nil {
		return alloc.ID, fmt.Errorf("Alloc. CreateAlloc failure %s", err.Error())
	}

	fmt.Println("-- hola 3 --")
	fmt.Println(validateErr)
	if validateErr != nil {
		return alloc.ID, validateErr
	}

	fmt.Println("-- hola 4")

	// Todo ha funcionado, enviar al scheduler
	err = s.scheduler.Apply(job, alloc)
	if err != nil {
		return alloc.ID, fmt.Errorf("Alloc. Scheduler apply failure %s", err.Error())
	}

	return alloc.ID, nil
}

func (s *Server) UpdateAllocStatus(req *structs.AllocStatusRequest) {
	fmt.Println("-- update alloc status --")
	fmt.Println(req)

	alloc := structs.Alloc{
		ID:     req.AllocID,
		Status: req.AllocStatus,
	}

	err := s.stateStore.UpdateAlloc(&alloc)
	if err != nil {
		fmt.Printf("CANNOT UPDATE VALUE CORRECTLY\n")
	}
}

func (s *Server) UpdateAllocDone(req *structs.AllocDoneRequest) {
	fmt.Println("-- update alloc done --")
	fmt.Println(req)

	var res map[string]interface{}

	fmt.Println("-- res --")
	fmt.Println(req)
	fmt.Println("-- json output --")
	fmt.Println(req.JSONOutput)

	outputErr := json.Unmarshal([]byte(req.JSONOutput), &res)

	// Adquirir el objeto alloc
	alloc, err := s.stateStore.AllocByID(req.AllocID)
	if err != nil {
		fmt.Println("FAILURE 1 " + err.Error())
		return
	}

	if outputErr != nil {
		alloc.Status = structs.FAILED
	}

	alloc.Output = res
	alloc.Logoutput = req.LogOutput
	alloc.Logerror = req.LogError

	// Adquirir el objeto job
	job, err := s.stateStore.JobByID(alloc.JobID)
	if err != nil {
		fmt.Println("FAILURE 2 " + err.Error())
		return
	}

	// Comprobar si el output es valido
	validateErr := ValidateOutput(s.config.OnedataPath, job, alloc)
	if validateErr != nil {
		fmt.Println("ERRR 1")
		fmt.Println(validateErr)

		alloc.Error = validateErr.Error()
		alloc.Status = structs.VALIDATIONOUTPUTFAILURE
	} else {
		alloc.Status = structs.COMPLETED
	}

	fmt.Println("VALIDO OUTPUT")

	fmt.Println("-- alloc para update --")
	fmt.Println(alloc)

	// Write to alloc with the newvalues
	err = s.stateStore.UpdateAlloc(alloc)
	if err != nil {
		fmt.Println("ERR 2")
		fmt.Println(err)
	}

	if s.checkIfPipelineApplies(alloc) {
		err = s.applyPipeline(alloc)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) applyPipeline(alloc *structs.Alloc) error {

	workflow, err := s.stateStore.WorkflowByID(alloc.WorkflowID)
	if err != nil {
		return err
	}

	// Sacar el workflow allocation
	workflowAlloc, err := s.stateStore.WorkflowAllocByID(alloc.WorkflowAllocID)
	if err != nil {
		return err
	}

	newStatus, err := s.applyPipelineImpl(workflow, workflowAlloc, alloc)
	workflowAlloc.Status = newStatus

	fmt.Println("-- el estado deberia de ser --")
	fmt.Println(structs.COMPLETED)
	fmt.Println(newStatus)
	fmt.Println(workflowAlloc)

	err = s.stateStore.UpdateWorkflowAlloc(workflowAlloc)
	if err != nil {
		fmt.Println("-- error --")
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Server) applyPipelineImpl(workflow *structs.Workflow, workflowAlloc *structs.WorkflowAlloc, alloc *structs.Alloc) (structs.AllocStatus, error) {

	fmt.Println("#########################")
	fmt.Println("-- aplicar el pipeline --")
	fmt.Println("#########################")

	fmt.Println("-- workflow --")
	fmt.Println(workflow)
	fmt.Println(alloc.WorkflowAllocID)

	fmt.Println("- workflow alloc -")
	fmt.Println(workflowAlloc)

	// Las definiciones de los trabajos
	jobs, err := s.stateStore.WorkflowJobs(workflow)
	if err != nil {
		return structs.FAILED, err
	}

	// Todos los allocations anteriores ya ejecutados
	allocs, err := s.stateStore.WorkflowAllocsByWorkflowAllocID(alloc.WorkflowAllocID)
	if err != nil {
		return structs.FAILED, err
	}

	// convert allocs to map with index is the job
	mapallocs := map[string]structs.Alloc{}
	for _, i := range *allocs {
		mapallocs[i.JobID] = i
	}

	// Detect the next step
	var position int
	for index, step := range workflow.Steps {
		if step.Job == alloc.JobID {
			position = index
			break
		}
	}

	// Not more jobs to run
	if position == len(workflow.Steps)-1 {
		fmt.Println("_ WORKFLOW ACABADO _")
		return structs.COMPLETED, err
	}

	nextStep := workflow.Steps[position+1]

	// Sacar el step en el que se encuentra ahora
	newalloc, err := CreateAllocFromStep(workflowAlloc, &nextStep, jobs, &workflow.Variables, &workflowAlloc.Input, &mapallocs)
	if err != nil {
		return structs.FAILED, err
	}

	fmt.Println("-- new allocation --")
	fmt.Println(newalloc)

	_, err = s.CreateAlloc(newalloc)
	if err != nil {
		return structs.FAILED, err
	}

	return structs.RUNNING, nil
}

func (s *Server) checkIfPipelineApplies(alloc *structs.Alloc) bool {
	return alloc.WorkflowID != ""
}
