package pergamo

import (
	"fmt"
	"pergamo/pergamo/structs"
	"testing"
)

func TestCompleteAlloc(t *testing.T) {

	job := structs.Job{
		ID:          "jobid",
		DriverImage: "sahkp/medisample1",
		Input: structs.JSON{
			"size": "number",
		},
		Output: structs.JSON{
			"image": "image.jpeg",
		},
	}

	alloc := structs.Alloc{
		JobID: "jobid",
		Input: structs.JSONGeneric{
			"size": 1,
		},
		Output: structs.JSONGeneric{},
	}

	completeAlloc(&job, &alloc)
	fmt.Println(alloc)
}
