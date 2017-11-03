package pergamo

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pergamo/pergamo/structs"
	"strconv"

	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type HttpResult struct {
	Result  interface{}
	HasNext bool
}

type HttpCreate struct {
	ID string
}

func SuccessMessage(c echo.Context, object interface{}) error {
	obj := map[string]interface{}{
		"Status": "OK",
		"Result": object,
	}

	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, string(data))
}

func ErrorMessage(c echo.Context, message string) error {
	fmt.Println("=> Error: " + message)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Status":  "BAD",
		"Message": message,
	})
}

func parseContextInt(c echo.Context, queryParam string, defaultValue int) (int, error) {
	if str := c.QueryParam(queryParam); str != "" {
		number, err := strconv.Atoi(str)
		if err != nil {
			return 0, nil
		}
		return number, nil
	}

	return defaultValue, nil
}

func parseOffset(c echo.Context) (*structs.QueryOffset, error) {

	page, err := parseContextInt(c, "page", 0)
	if err != nil {
		return nil, err
	}

	offset, err := parseContextInt(c, "offset", 10)
	if err != nil {
		return nil, err
	}

	return &structs.QueryOffset{
		Page:   page,
		Offset: offset,
	}, nil
}

type HttpServer struct {
	config *Config
	HTTP   *echo.Echo
	Server *Server
}

func (ht *HttpServer) listen() {
	ht.HTTP.Start(":" + strconv.Itoa(ht.config.APIPort))
}

func (ht *HttpServer) registerEndpoints() {
	ht.HTTP.GET("/", ht.GetIndex)

	// Images
	ht.HTTP.GET("/images", ht.GetImages)
	ht.HTTP.GET("/images/:imageid/download", ht.GetImagesImageidDownload)
	ht.HTTP.POST("/images", ht.PostImages)

	// Job definitions
	ht.HTTP.GET("/jobs", ht.GetJobs)
	ht.HTTP.POST("/jobs", ht.PostJobs)
	ht.HTTP.GET("/jobs/:jobid", ht.GetJobsJobid)
	ht.HTTP.DELETE("/jobs/:jobid", ht.DeleteJobsJobid)

	// Job allocations
	ht.HTTP.GET("/allocs", ht.GetAllocs)
	ht.HTTP.POST("/allocs", ht.PostAllocs)
	ht.HTTP.GET("/allocs/:allocid", ht.GetAllocsAllocid)
	ht.HTTP.GET("/allocs/:allocid/download", ht.GetAllocsAllocidDownload)
	ht.HTTP.POST("/allocs/:allocid/running", ht.PostAllocsAllocIDRunning)
	ht.HTTP.POST("/allocs/:allocid/done", ht.PostAllocsAllocIDDone)

	// Workflow definitions
	ht.HTTP.GET("/workflows", ht.GetWorkflows)
	ht.HTTP.GET("/workflows/:workflowid", ht.GetWorkflowsWorkflowid)
	ht.HTTP.POST("/workflows", ht.PostWorkflows)

	// Workflow allocations
	ht.HTTP.GET("/workallocs", ht.GetWorkallocs)
	ht.HTTP.GET("/workallocs/:workallocid", ht.GetWorkallocsWorkallocid)
	ht.HTTP.POST("/workallocs", ht.PostWorkallocs)
}

func (ht *HttpServer) PostWorkallocs(c echo.Context) error {
	id := c.FormValue("id")
	if id == "" {
		return ErrorMessage(c, "job not provided")
	}

	attrs := c.FormValue("attrs")
	if attrs == "" {
		return ErrorMessage(c, "attrs not provided")
	}

	var attrsJSON structs.JSONGeneric
	err := json.Unmarshal([]byte(attrs), &attrsJSON)
	if err != nil {
		return ErrorMessage(c, "Cannot parse json generic on allocs call "+err.Error())
	}

	workflowAlloc := structs.WorkflowAlloc{
		WorkflowID: id,
		Input:      attrsJSON,
	}

	fmt.Println("workflow alloc")
	fmt.Println(workflowAlloc)

	err = ht.Server.ApplyWorkflowAlloc(&workflowAlloc)
	if err != nil {
		return ErrorMessage(c, "Failed to create job "+err.Error())
	}

	return SuccessMessage(c, HttpCreate{ID: workflowAlloc.ID})
}

func (ht *HttpServer) PostWorkflows(c echo.Context) error {
	content := c.FormValue("content")
	if content == "" {
		return ErrorMessage(c, "Not content found")
	}

	//structs.Workflow
	//unmarhsall

	var workflow structs.Workflow
	err := json.Unmarshal([]byte(content), &workflow)
	if err != nil {
		panic(err)
	}

	req := structs.WorkflowPostRequest{Workflow: &workflow}
	var resp structs.WorkflowPostResponse

	err = ht.Server.PostWorkflow(&req, &resp)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	fmt.Println("-- workflow creado --")
	fmt.Println(workflow)

	return SuccessMessage(c, resp.ID)
}

func (ht *HttpServer) GetImagesImageidDownload(c echo.Context) error {
	imageid := c.Param("imageid")
	if imageid == "" {
		return ErrorMessage(c, "imageid not found")
	}

	path, name, err := ht.Server.CheckDownloadImage(imageid)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return c.Attachment(path, name)
}

func (ht *HttpServer) PostAllocs(c echo.Context) error {

	fmt.Println("-- POST ALLOCS --")
	fmt.Println(c.FormParams())

	job := c.FormValue("job")
	if job == "" {
		return ErrorMessage(c, "job not provided")
	}

	attrs := c.FormValue("attrs")
	if attrs == "" {
		return ErrorMessage(c, "attrs not provided")
	}

	var attrsJSON structs.JSONGeneric
	err := json.Unmarshal([]byte(attrs), &attrsJSON)
	if err != nil {
		return ErrorMessage(c, "Cannot parse json generic on allocs call "+err.Error())
	}

	alloc := structs.Alloc{
		JobID:  job,
		Input:  attrsJSON,
		Output: structs.JSONGeneric{},
	}

	req := structs.AllocPostRequest{Alloc: alloc}
	var resp structs.AllocPostResponse

	err = ht.Server.PostAlloc(&req, &resp)
	if err != nil {
		return ErrorMessage(c, "Failed to create job "+err.Error())
	}

	return SuccessMessage(c, HttpCreate{ID: resp.ID})
}

func stringToBool(s string) (bool, error) {
	s = strings.ToLower(s)

	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	return false, fmt.Errorf("String to bool failed %s", s)
}

func (ht *HttpServer) PostJobs(c echo.Context) error {

	fmt.Println("-- POST JOBS --")
	//fmt.Println(c.FormParams())

	// name
	name := c.FormValue("name")
	if name == "" {
		return ErrorMessage(c, "name not provided")
	}

	// image
	image := c.FormValue("image")
	if image == "" {
		return ErrorMessage(c, "image not provided")
	}

	// attrs
	input := c.FormValue("input")
	if input == "" {
		return ErrorMessage(c, "input not provided")
	}

	output := c.FormValue("output")
	if output == "" {
		return ErrorMessage(c, "input not provided")
	}

	checkfilesStr := c.FormValue("checkfiles")
	if checkfilesStr == "" {
		return ErrorMessage(c, "checkfiles not provided")
	}

	checkfiles, err := stringToBool(checkfilesStr)
	if err != nil {
		return ErrorMessage(c, fmt.Sprintf("Checkfiles: %v", err))
	}

	// file
	file, err := c.FormFile("file")
	if err != nil {
		if !strings.Contains(err.Error(), "no such file") {
			return ErrorMessage(c, "file upload header error "+err.Error())
		}
	}

	//fmt.Println(file)

	var inputJSON structs.JSON
	err = json.Unmarshal([]byte(input), &inputJSON)
	if err != nil {
		return ErrorMessage(c, "Cannot parse input")
	}

	var outputJSON structs.JSON
	err = json.Unmarshal([]byte(output), &outputJSON)
	if err != nil {
		return ErrorMessage(c, "Cannot parse output")
	}

	fmt.Println("-- creacion --")
	fmt.Println(inputJSON)
	fmt.Println(outputJSON)

	job := structs.Job{
		ID:          name,
		DriverImage: image,
		Input:       inputJSON,
		Output:      outputJSON,
		Checkfiles:  checkfiles,
	}

	req := structs.JobPostRequest{Job: job, File: file}
	var resp structs.JobPostResponse

	err = ht.Server.PostJob(&req, &resp)
	if err != nil {
		return ErrorMessage(c, "Failed to create job "+err.Error())
	}

	return SuccessMessage(c, HttpCreate{ID: resp.ID})
}

func (ht *HttpServer) GetIndex(c echo.Context) error {
	return c.String(http.StatusOK, "index")
}

func (ht *HttpServer) GetAllocsAllocidDownload(c echo.Context) error {

	fmt.Println("DOWNLOAD")

	allocid := c.Param("allocid")
	if allocid == "" {
		return ErrorMessage(c, "allocid not found")
	}

	c.Response().Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", allocid))

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	zw := zip.NewWriter(writer)

	req := structs.AllocDownloadRequest{Allocid: allocid, Writer: zw}
	var resp structs.AllocDownloadResponse

	err := ht.Server.Download(&req, &resp)
	if err != nil {
		return ErrorMessage(c, "Failed to download image "+err.Error())
	}

	err = zw.Close()
	if err != nil {
		ErrorMessage(c, "Failed to close "+err.Error())
	}

	//return c.File("/tmp/some.json")
	return c.Blob(http.StatusOK, "application/zip", b.Bytes())
}

func (ht *HttpServer) PostImages(c echo.Context) error {
	fmt.Println("-- values --")
	fmt.Println(c.FormParams())

	imagename := c.FormValue("name")
	if imagename == "" {
		return ErrorMessage(c, "name not provided")
	}

	format := c.FormValue("format")
	if format == "" {
		return ErrorMessage(c, "format not provided")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return ErrorMessage(c, "Could not load the file "+err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return ErrorMessage(c, "Could not read the file "+err.Error())
	}
	defer src.Close()

	content, err := ioutil.ReadAll(src)
	if err != nil {
		return ErrorMessage(c, "Failed to read content "+err.Error())
	}

	req := structs.ImagePostRequest{Name: imagename, Format: format, Content: content}
	var resp structs.ImagePostResponse

	err = ht.Server.PostImage(&req, &resp)
	if err != nil {
		return ErrorMessage(c, "Failed to upload image "+err.Error())
	}

	return SuccessMessage(c, HttpCreate{ID: resp.ID})
}

func (ht *HttpServer) GetImages(c echo.Context) error {
	queryOffset, err := parseOffset(c)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	req := structs.ImageListRequest{QueryOffset: *queryOffset}
	var resp structs.ImageListResponse

	err = ht.Server.ImageList(&req, &resp)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, HttpResult{HasNext: resp.HasNext, Result: resp.Images})
}

func (ht *HttpServer) GetJobs(c echo.Context) error {
	queryOffset, err := parseOffset(c)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	req := structs.JobListRequest{QueryOffset: *queryOffset}
	var resp structs.JobListResponse

	err = ht.Server.JobList(&req, &resp)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, HttpResult{HasNext: resp.HasNext, Result: resp.Jobs})
}

func (ht *HttpServer) GetAllocs(c echo.Context) error {
	queryOffset, err := parseOffset(c)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	req := structs.AllocListRequest{QueryOffset: *queryOffset}
	var resp structs.AllocListResponse

	err = ht.Server.AllocList(&req, &resp)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, HttpResult{HasNext: resp.HasNext, Result: resp.Allocs})
}

func (ht *HttpServer) GetWorkflows(c echo.Context) error {

	queryOffset, err := parseOffset(c)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	req := structs.WorkflowListRequest{QueryOffset: *queryOffset}
	var resp structs.WorkflowListResponse

	err = ht.Server.WorkflowList(&req, &resp)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, HttpResult{HasNext: resp.HasNext, Result: resp.Workflows})
}

func (ht *HttpServer) GetWorkallocs(c echo.Context) error {
	queryOffset, err := parseOffset(c)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	req := structs.WorkflowAllocListRequest{QueryOffset: *queryOffset}
	var resp structs.WorkflowAllocListResponse

	err = ht.Server.WorkflowAllocList(&req, &resp)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, HttpResult{HasNext: resp.HasNext, Result: resp.WorkflowAllocs})
}

func (ht *HttpServer) GetJobsJobid(c echo.Context) error {
	jobid := c.Param("jobid")
	if jobid == "" {
		return ErrorMessage(c, "Jobid is not found in param query")
	}

	job, err := ht.Server.stateStore.JobByID(jobid)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, job)
}

func (ht *HttpServer) DeleteJobsJobid(c echo.Context) error {
	jobid := c.Param("jobid")
	if jobid == "" {
		return ErrorMessage(c, "Jobid is not found in param query")
	}

	err := ht.Server.stateStore.DeleteJobByID(jobid)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, jobid)
}

func (ht *HttpServer) GetWorkflowsWorkflowid(c echo.Context) error {
	workflowid := c.Param("workflowid")
	if workflowid == "" {
		return ErrorMessage(c, "workflowid is not found in param query")
	}

	workflow, err := ht.Server.stateStore.WorkflowByID(workflowid)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, workflow)
}

func (ht *HttpServer) GetWorkallocsWorkallocid(c echo.Context) error {
	workallocid := c.Param("workallocid")
	if workallocid == "" {
		return ErrorMessage(c, "workallocid is not found in param query")
	}

	workflow, err := ht.Server.stateStore.WorkflowAllocByID(workallocid)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, workflow)
}

func (ht *HttpServer) GetAllocsAllocid(c echo.Context) error {
	allocid := c.Param("allocid")
	if allocid == "" {
		return ErrorMessage(c, "Allocid is not found in param query")
	}

	req := structs.AllocSingleRequest{ID: allocid}
	var resp structs.AllocSingleResponse

	err := ht.Server.AllocByID(&req, &resp)
	if err != nil {
		return ErrorMessage(c, err.Error())
	}

	return SuccessMessage(c, resp.Alloc)
}

func (ht *HttpServer) PostAllocsAllocIDRunning(c echo.Context) error {
	allocid := c.Param("allocid")
	if allocid == "" {
		return ErrorMessage(c, "Allocid is not found in param query")
	}

	statusRequest := structs.AllocStatusRequest{
		AllocID:     allocid,
		AllocStatus: structs.RUNNING,
	}

	ht.Server.UpdateAllocStatus(&statusRequest)
	return c.String(http.StatusOK, "hola")
}

func (ht *HttpServer) PostAllocsAllocIDDone(c echo.Context) error {
	allocid := c.Param("allocid")
	if allocid == "" {
		return ErrorMessage(c, "Allocid is not found in param query")
	}

	logoutput := c.FormValue("logoutput")
	if logoutput == "" {
		// No hay logoutput
	}

	logerror := c.FormValue("logerror")
	if logerror == "" {
		// No hay logerror
	}

	jsonoutput := c.FormValue("jsonoutput")
	if jsonoutput == "" {
		// No hay json output
	}

	doneRequest := structs.AllocDoneRequest{
		AllocID:    allocid,
		LogOutput:  logoutput,
		LogError:   logerror,
		JSONOutput: jsonoutput,
	}

	ht.Server.UpdateAllocDone(&doneRequest)
	return c.String(http.StatusOK, "hola")
}

func NewHttpServer(server *Server, config *Config) (*HttpServer, error) {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	http := &HttpServer{
		config: config,
		HTTP:   e,
		Server: server,
	}

	http.registerEndpoints()

	return http, nil
}
