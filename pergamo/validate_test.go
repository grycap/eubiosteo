package pergamo

import (
	"fmt"
	"pergamo/pergamo/structs"
	"testing"
)

func TestAllocation(t *testing.T) {

	job := structs.Job{
		ID: "a",
		Input: structs.JSON{
			"size": "slice.number",
		},
	}

	alloc := structs.Alloc{
		Input: structs.JSONGeneric{
			"size": "[1]",
		},
	}

	fmt.Println(job)
	fmt.Println(alloc)

	err := ValidateInput("/", &job, &alloc)

	fmt.Println(err)
}

func TestWorkflowExpect(t *testing.T) {
}
