package structs

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewJob(t *testing.T) {

	var jobFixture = `
{
	"id": "hola",
	"input": {
		"a": 444
	}
}
`

	var job Job
	err := json.Unmarshal([]byte(jobFixture), &job)
	if err != nil {
		panic(err)
	}

	fmt.Println(job)
}
