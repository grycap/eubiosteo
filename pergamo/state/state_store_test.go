package state

import (
	"fmt"
	"pergamo/pergamo/structs"
	"testing"
)

var purgeDB = `
PRAGMA writable_schema = 1;
delete from sqlite_master where type in ('table', 'index', 'trigger');
PRAGMA writable_schema = 0;
`

func Initialize(t *testing.T, initialize bool) *StateStore {
	state, err := NewStateStore("", "")
	if err != nil {
		t.Fatal(err)
	}

	if initialize {
		err = state.InitializeSchema()
		if err != nil {
			t.Fatalf("Fallo al crear el esquema de prueba %s", err.Error())
		}
	}

	return state
}

func PurgeSchema(state *StateStore, t *testing.T) {
	_, err := state.db.Exec(purgeDB)
	if err != nil {
		t.Fatalf("No se ha podido eliminar la bd en purge. Error: %s", err.Error())
	}
}

func TestCreateSchema(t *testing.T) {
	state := Initialize(t, true)
	defer PurgeSchema(state, t)
}

func TestCreateJob(t *testing.T) {
	state := Initialize(t, true)
	defer PurgeSchema(state, t)

	newjob := structs.Job{
		ID:          "main",
		DriverImage: "image",
		Zip:         "zip",
		Input: structs.JSON{
			"input1": "number",
		},
	}

	jobid, err := state.CreateJob(&newjob)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(jobid)

	job, err := state.JobByID(jobid)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(job)

	newalloc := structs.Alloc{
		JobID: jobid,
		Input: structs.JSONGeneric{
			"input1": 1,
		},
	}

	err = state.CreateAlloc(&newalloc)
	if err != nil {
		panic(err)
	}

	fmt.Println(newalloc.ID)

	alloc, err := state.AllocByID(newalloc.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(alloc)

	// update alloc

	newalloc.Initialtime = makeTimestamp()

	err = state.UpdateAlloc(&newalloc)
	if err != nil {
		panic(err)
	}

	alloc, err = state.AllocByID(newalloc.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("-- nuevo recibido --")
	fmt.Println(alloc)
}
