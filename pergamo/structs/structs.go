package structs

import (
	"archive/zip"
	"fmt"
	"mime/multipart"
	"strconv"
)

type JSON map[string]string

type Job struct {
	ID     string `db:"id"`
	Input  JSON   `db:"input"`
	Output JSON   `db:"output"`

	Checkfiles bool `db:"checkfiles"`

	DriverImage string `db:"driverimage"`
	Zip         string `db:"zip"`
}

type Image struct {
	ID string `db:"id"`

	Name string `db:"name"`
	Path string `db:"path"`

	Format string `db:"format"`
}

type AllocStatus int

const (
	PENDING AllocStatus = iota + 1
	RUNNING
	COMPLETED
	CANCELLED
	FAILED
	UNKNOWN
	VALIDATIONINPUTFAILURE
	VALIDATIONOUTPUTFAILURE
)

func (a *AllocStatus) String() string {
	switch *a {
	case PENDING:
		return "Pending"
	case RUNNING:
		return "Running"
	case COMPLETED:
		return "Completed"
	case CANCELLED:
		return "Cancelled"
	case FAILED:
		return "Failed"
	case UNKNOWN:
		return "Unknown"
	case VALIDATIONINPUTFAILURE:
		return "ValidationInputFailure"
	case VALIDATIONOUTPUTFAILURE:
		return "ValidationOutputFailure"
	}

	panic(fmt.Errorf("String not found for %d", *a))
}

type JSONGeneric map[string]interface{}

type Timestamp int64

type Alloc struct {
	ID string `db:"id"`

	WorkflowID      string `db:"workflowid"`
	WorkflowAllocID string `db:"workflowallocid"`
	JobID           string `db:"jobid"`

	Input  JSONGeneric `db:"input"`
	Output JSONGeneric `db:"output"`

	Status     AllocStatus `db:"status"`
	ScriptPath string      `db:"scriptpath"`

	Logoutput string `db:"logoutput"`
	Logerror  string `db:"logerror"`

	Error string `db:"error"`

	Initialtime Timestamp `db:"initialtime"`
	Finaltime   Timestamp `db:"finaltime"`
	Elapsedtime Timestamp `db:"elapsedtime"`
}

// -- WORKFLOW --

type Step struct {
	Job   string
	Param JSONGeneric
}

type Steps []Step

type Workflow struct {
	ID        string      `db:"id"`
	Variables JSONGeneric `db:"variables"`
	Steps     Steps       `db:"steps"`
	Entry     JSON        `db:"entry"`
}

type WorkflowAlloc struct {
	ID         string `db:"id"`
	WorkflowID string `db:"workflowid"`

	Status AllocStatus `db:"status"`

	Input  JSONGeneric `db:"input"`
	Output JSONGeneric `db:"output"`

	Initialtime int64 `db:"initialtime"`
	Finaltime   int64 `db:"finaltime"`
	Elapsedtime int64 `db:"elapsedtime"`

	Allocs []Alloc
}

// TODO. Querycursor ResponseCursor to allow cursor pagination

type WriteMeta struct {
	ID string
}

type QueryOffset struct {
	Page   int
	Offset int
}

func (q *QueryOffset) GoString(checknext bool) string {
	offset := q.Offset
	if checknext {
		offset++
	}

	return " LIMIT " + strconv.Itoa(offset) + " OFFSET " + strconv.Itoa(q.Page*q.Offset)
}

type ResponseOffset struct {
	HasNext bool
}

type AllocDoneRequest struct {
	AllocID    string
	LogOutput  string
	LogError   string
	JSONOutput string

	Err error
}

type AllocStatusRequest struct {
	AllocID     string
	AllocStatus AllocStatus
}

type AllocSingleRequest struct {
	ID string
}

type AllocSingleResponse struct {
	Alloc *Alloc
}

// LIST REQUESTS

type ImageListRequest struct {
	QueryOffset
}

type ImageListResponse struct {
	Images []Image
	ResponseOffset
}

type JobListRequest struct {
	QueryOffset
}

type JobListResponse struct {
	Jobs []Job
	ResponseOffset
}

type AllocListRequest struct {
	QueryOffset
}

type AllocListResponse struct {
	Allocs []Alloc
	ResponseOffset
}

type WorkflowListRequest struct {
	QueryOffset
}

type WorkflowListResponse struct {
	Workflows []Workflow
	ResponseOffset
}

type WorkflowAllocListRequest struct {
	QueryOffset
}

type WorkflowAllocListResponse struct {
	WorkflowAllocs []WorkflowAlloc
	ResponseOffset
}

type ImagePostRequest struct {
	Name    string
	Format  string
	Content []byte
}

type ImagePostResponse struct {
	WriteMeta
}

type JobPostRequest struct {
	Job  Job
	File *multipart.FileHeader
}

type JobPostResponse struct {
	WriteMeta
}

type AllocPostRequest struct {
	Alloc Alloc
}

type AllocPostResponse struct {
	WriteMeta
}

type AllocDownloadRequest struct {
	Allocid string
	Writer  *zip.Writer
}

type AllocDownloadResponse struct {
}

type WorkflowPostRequest struct {
	Workflow *Workflow
}

type WorkflowPostResponse struct {
	WriteMeta
}

type WorkflowAllocPostRequest struct {
	WorkflowAlloc WorkflowAlloc
}

type WorkflowAllocPostResponse struct {
	WriteMeta
}
