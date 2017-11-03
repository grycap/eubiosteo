package state

import (
	"fmt"
	"io/ioutil"
	"os"
	"pergamo/pergamo/structs"
	"time"

	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type StateStore struct {
	db          *sqlx.DB
	onedataPath string
	dbdriver    string
	dbhost      string
}

func NewStateStore(onedataPath, dbdriver, dbhost string) (*StateStore, error) {

	db, err := sqlx.Open(dbdriver, dbhost)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(onedataPath); os.IsNotExist(err) {
		// Create storage folder if does not exists. It should happen with local deployments.
		err = os.MkdirAll(onedataPath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("Error creating storage folder %s: %s", onedataPath, err.Error())
		}
	}

	state := &StateStore{
		db:          db,
		onedataPath: onedataPath,
		dbdriver:    dbdriver,
		dbhost:      dbhost,
	}

	if dbdriver == "sqlite3" {
		if dbhost != ":memory:" {
			if _, err := os.Stat(dbhost); err == nil {
				goto DONE
			}
		}

		fmt.Println("Initialize schema")
		err := state.InitializeSchema()
		if err != nil {
			return nil, err
		}
	}

DONE:
	return state, nil
}

func (s *StateStore) InitializeSchema() error {
	_, err := s.db.Exec(Schema)
	return err
}

func (s *StateStore) SaveImageToStorage(image *structs.Image, content []byte) error {
	return ioutil.WriteFile(filepath.Join(s.onedataPath, image.Path), content, 0644)
}

func (s *StateStore) SaveImage(image *structs.Image) error {
	imageid := structs.GenerateUUID()

	image.ID = imageid
	image.Path = imageid

	_, err := s.db.NamedExec("INSERT INTO image (id, name, path, format) VALUES (:id, :name, :path, :format)", image)
	return err
}

func (s *StateStore) DeleteJobByID(jobid string) error {
	_, err := s.db.Exec("DELETE FROM job WHERE id=$1", jobid)
	return err
}

func (s *StateStore) GetImage(id string) (*structs.Image, error) {
	var image structs.Image
	err := s.db.Get(&image, "SELECT * FROM image WHERE ID=$1", id)

	return &image, err
}

func (s *StateStore) GetImages(query *structs.QueryOffset) (*[]structs.Image, bool, error) {
	var images []structs.Image
	next := false

	queryStr := "SELECT * FROM image"
	if offsetQuery := query.GoString(true); offsetQuery != "" {
		queryStr += offsetQuery
	}

	err := s.db.Select(&images, queryStr)
	if len(images) == query.Offset+1 {
		next = true
		images = images[0:query.Offset]
	}

	if len(images) == 0 {
		images = []structs.Image{}
	}

	return &images, next, err
}

func (s *StateStore) GetJobs(query *structs.QueryOffset) (*[]structs.Job, bool, error) {
	var jobs []structs.Job
	next := false

	queryStr := "SELECT * FROM job"
	if offsetQuery := query.GoString(true); offsetQuery != "" {
		queryStr += offsetQuery
	}

	err := s.db.Select(&jobs, queryStr)
	if len(jobs) == query.Offset+1 {
		next = true
		jobs = jobs[0:query.Offset]
	}

	if len(jobs) == 0 {
		jobs = []structs.Job{}
	}
	return &jobs, next, err
}

func (s *StateStore) GetAllocs(query *structs.QueryOffset) (*[]structs.Alloc, bool, error) {
	var allocs []structs.Alloc
	next := false

	queryStr := "SELECT * FROM alloc"
	if offsetQuery := query.GoString(true); offsetQuery != "" {
		queryStr += offsetQuery
	}

	err := s.db.Select(&allocs, queryStr)
	if len(allocs) == query.Offset+1 {
		next = true
		allocs = allocs[0:query.Offset]
	}

	if len(allocs) == 0 {
		allocs = []structs.Alloc{}
	}
	return &allocs, next, err
}

func (s *StateStore) GetWorkflows(query *structs.QueryOffset) (*[]structs.Workflow, bool, error) {
	var workflows []structs.Workflow
	next := false

	queryStr := "SELECT * FROM workflow"
	if offsetQuery := query.GoString(true); offsetQuery != "" {
		queryStr += offsetQuery
	}

	err := s.db.Select(&workflows, queryStr)
	if len(workflows) == query.Offset+1 {
		next = true
		workflows = workflows[0:query.Offset]
	}

	if len(workflows) == 0 {
		workflows = []structs.Workflow{}
	}
	return &workflows, next, err
}

func (s *StateStore) GetWorkflowAllocs(query *structs.QueryOffset) (*[]structs.WorkflowAlloc, bool, error) {
	var workflowallocs []structs.WorkflowAlloc
	next := false

	queryStr := "SELECT * FROM workflowalloc"
	if offsetQuery := query.GoString(true); offsetQuery != "" {
		queryStr += offsetQuery
	}

	err := s.db.Select(&workflowallocs, queryStr)
	if len(workflowallocs) == query.Offset+1 {
		next = true
		workflowallocs = workflowallocs[0:query.Offset]
	}

	if len(workflowallocs) == 0 {
		workflowallocs = []structs.WorkflowAlloc{}
	}
	return &workflowallocs, next, err
}

func (s *StateStore) CreateJob(job *structs.Job) error {
	_, err := s.db.NamedExec("INSERT INTO job (id, input, output, driverimage, zip, checkfiles) VALUES (:id, :input, :output, :driverimage, :zip, :checkfiles)", job)
	return err
}

func (s *StateStore) CreateWorkflow(workflow *structs.Workflow) error {
	_, err := s.db.NamedExec("INSERT INTO workflow (id, variables, steps, entry) VALUES (:id, :variables, :steps, :entry)", workflow)
	return err
}

func jobsSliceToMap(jobs *[]structs.Job) *map[string]structs.Job {
	mapjobs := map[string]structs.Job{}
	for _, i := range *jobs {
		mapjobs[i.ID] = i
	}
	return &mapjobs
}

func (s *StateStore) CreateWorkflowAlloc(workflowalloc *structs.WorkflowAlloc) error {
	id := structs.GenerateUUID()
	workflowalloc.ID = id

	_, err := s.db.NamedExec("INSERT INTO workflowalloc (id, workflowid, status, input, output, initialtime, finaltime, elapsedtime) VALUES (:id, :workflowid, :status, :input, :output, :initialtime, :finaltime, :elapsedtime)", workflowalloc)
	return err
}

func (s *StateStore) UpdateWorkflowAlloc(workflowalloc *structs.WorkflowAlloc) error {

	fmt.Println("-- status ---")
	fmt.Println("- otro -")
	fmt.Println(workflowalloc.Status)

	_, err := s.db.NamedExec("UPDATE workflowalloc SET status=:status WHERE id=:id", workflowalloc)
	return err
}

func (s *StateStore) AllocsByWorkflowID(id string) (*map[string]structs.Alloc, error) {
	// todos los allocs del workflow por su id
	var allocs []structs.Alloc
	err := s.db.Select(&allocs, "SELECT * FROM alloc WHERE workflowallocid=$1", id)

	mapallocs := map[string]structs.Alloc{}
	for _, i := range allocs {
		mapallocs[i.JobID] = i
	}

	return &mapallocs, err
}

func (s *StateStore) WorkflowJobs(workflow *structs.Workflow) (*map[string]structs.Job, error) {
	ids := []string{}
	for _, i := range workflow.Steps {
		ids = append(ids, i.Job)
	}

	jobs, err := s.JobsByIDs(ids)
	if err != nil {
		return nil, err
	}

	mapjobs := map[string]structs.Job{}
	for _, i := range *jobs {
		mapjobs[i.ID] = i
	}

	return &mapjobs, err
}

func (s *StateStore) WorkflowByID(id string) (*structs.Workflow, error) {
	var workflow structs.Workflow
	err := s.db.Get(&workflow, "SELECT * FROM workflow WHERE id=$1", id)

	return &workflow, err
}

func (s *StateStore) JobsByIDs(ids []string) (*[]structs.Job, error) {
	jobs := []structs.Job{}
	for _, i := range ids {
		job, err := s.JobByID(i)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, *job)
	}

	return &jobs, nil
}

func (s *StateStore) JobByID(id string) (*structs.Job, error) {
	var job structs.Job
	err := s.db.Get(&job, "SELECT * FROM job WHERE id=$1", id)

	return &job, err
}

func makeTimestamp() structs.Timestamp {
	return structs.Timestamp(time.Now().UnixNano() / int64(time.Millisecond))
}

func (s *StateStore) WorkflowAllocsByWorkflowAllocID(workflowallocid string) (*[]structs.Alloc, error) {
	var allocs []structs.Alloc
	err := s.db.Select(&allocs, "SELECT * FROM alloc WHERE workflowallocid=$1", workflowallocid)

	return &allocs, err
}

func (s *StateStore) WorkflowAllocByID(id string) (*structs.WorkflowAlloc, error) {
	var workflowalloc structs.WorkflowAlloc
	fmt.Println("-a")
	err := s.db.Get(&workflowalloc, "SELECT * FROM workflowalloc WHERE id=$1", id)

	// Descargar tambien los otros allocs
	allocs, err := s.WorkflowAllocsByWorkflowAllocID(id)
	if err != nil {
		return nil, err
	}

	fmt.Println("-- allocs --")
	fmt.Println(allocs)

	workflowalloc.Allocs = *allocs

	fmt.Println("-b")
	return &workflowalloc, err
}

// CreateAlloc envia dentor de alloc los datos extra de la ejecucion
func (s *StateStore) CreateAlloc(alloc *structs.Alloc) error {

	allocid := structs.GenerateUUID()
	//allocpathid := structs.GenerateUUID()

	alloc.ID = allocid
	alloc.ScriptPath = allocid
	alloc.Initialtime = makeTimestamp()

	_, err := s.db.NamedExec("INSERT INTO alloc (id, workflowid, workflowallocid, jobid, input, output, status, scriptpath, logoutput, logerror, initialtime, error) VALUES (:id, :workflowid, :workflowallocid, :jobid, :input, :output, :status, :scriptpath, :logoutput, :logerror, :initialtime, :error)", alloc)
	return err
}

func (s *StateStore) AllocByID(id string) (*structs.Alloc, error) {
	var alloc structs.Alloc
	err := s.db.Get(&alloc, "SELECT * FROM alloc WHERE id=$1", id)

	return &alloc, err
}

func (s *StateStore) UpdateAlloc(alloc *structs.Alloc) error {

	lasttime := makeTimestamp()
	elapsed := lasttime - alloc.Initialtime

	alloc.Finaltime = lasttime
	alloc.Elapsedtime = elapsed

	fmt.Println("-- final alloc to update --")
	fmt.Println(alloc)
	fmt.Println(alloc.Status)

	_, err := s.db.NamedExec("UPDATE alloc SET output=:output, status=:status, error=:error, logoutput=:logoutput, logerror=:logerror, finaltime=:finaltime, elapsedtime=:elapsedtime WHERE id=:id", alloc)
	//_, err := s.db.Exec("UPDATE alloc SET scriptpath=$1 WHERE id=$2", "XXXXX", alloc.ID)
	return err
}
