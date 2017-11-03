package script

import (
	"encoding/json"
	"fmt"
	"pergamo/pergamo/structs"
	"pergamo/pergamo/types"
)

type Script interface {
	Plan(onedata string, job *structs.Job, alloc *structs.Alloc, addrs string) (string, error)
}

type ScriptState struct {
	JobID       string
	DockerImage string
	Input       string
}

func createInputJson(job *structs.Job, alloc *structs.Alloc) (string, error) {

	res := map[string]interface{}{}
	for name, event := range job.Input {

		inputEvent, ok := alloc.Input[name]
		if !ok {
			return "", fmt.Errorf("ValidateAlloc. El input de nombre %s no esta", name)
		}

		jobType, err := types.NewType(event)
		if err != nil {
			return "", fmt.Errorf("ValidateAlloc. El input de nombre %s no tiene el bien de codigo %s", name, event)
		}

		if types.IsFileType(name) {
			// pasar solo la referencia del archivo
			inputEventStr, ok := inputEvent.(string)
			if !ok {
				return "", fmt.Errorf("ValidateAlloc. Cannot cast filetype to string")
			}

			res[name] = inputEventStr
		} else {

			// hacer parse y pasar la version parseada
			err = jobType.Parse(inputEvent)
			if err != nil {
				return "", fmt.Errorf("ValidateAlloc. Error al validar input con parametro para %s con error %s", name, err.Error())
			}

			res[name] = jobType.Value()
		}

	}

	data, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("Cannot convert to marshall %s", err.Error())
	}

	return string(data), nil
}
