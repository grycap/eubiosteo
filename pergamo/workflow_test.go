package pergamo

import (
	"fmt"
	"pergamo/pergamo/structs"
	"testing"
)

func TestParamSource(t *testing.T) {

	fmt.Println(parseParamSource("@step_job0_output_b"))
	fmt.Println(parseParamSource("@variable_s"))
}

func TestWorkflow(t *testing.T) {

	job0 := structs.Job{
		ID: "job0",
		Input: structs.JSON{
			"a":  "number",
			"a1": "slice.number",
		},
		Output: structs.JSON{
			"b": "slice.number",
		},
	}

	job1 := structs.Job{
		ID: "job1",
		Input: structs.JSON{
			"c": "slice.number",
		},
		Output: structs.JSON{
			"d": "number",
		},
	}

	workflow := structs.Workflow{
		Variables: structs.JSONGeneric{
			"entry": []int{1, 2, 3},
		},
		Steps: []structs.Step{
			structs.Step{
				Job: "job0",
				Param: structs.JSONGeneric{
					"a":  1,
					"a1": "@variable_entry",
				},
			},
			structs.Step{
				Job: "job1",
				Param: structs.JSONGeneric{
					"c": "@step_job0_output_b",
				},
			},
		},
	}

	jobs := []structs.Job{
		job0,
		job1,
	}

	fmt.Println(ValidateWorkflow(&workflow, &jobs))

}
