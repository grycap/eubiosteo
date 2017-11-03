package pergamo

import (
	"fmt"
	"pergamo/pergamo/structs"
	"pergamo/pergamo/types"
	"regexp"
)

func jobsToMap(jobs []structs.Job) map[string]structs.Job {
	mapjobs := map[string]structs.Job{}
	for _, job := range jobs {
		mapjobs[job.ID] = job
	}
	return mapjobs
}

type paramSourceType int

const (
	sourceVariable paramSourceType = iota + 1
	sourceEvent
	sourceInput
)

// paramsource de tipo @variable o @event.job.output.param
type paramSource struct {
	Type  paramSourceType
	Name  string // Nombre de la variable o nombre del job
	Loc   string // input | output
	Param string // El parametro que se quiere sacar del job
}

func isGoodLocation(s string) bool {
	return s == "input" || s == "output"
}

func parseParamSource(s interface{}) (*paramSource, bool, error) {
	s0, ok := s.(string)
	if !ok {
		return nil, false, nil
	}

	var source *paramSource

	inputMatch := "@input_(\\w+)"

	match := regexp.MustCompile(inputMatch).FindStringSubmatch(s0)
	if len(match) == 2 {
		source = &paramSource{
			Type: sourceInput,
			Name: match[1],
		}
	}

	variableMatch := "@variable_(\\w+)"

	match = regexp.MustCompile(variableMatch).FindStringSubmatch(s0)
	if len(match) == 2 {
		source = &paramSource{
			Type: sourceVariable,
			Name: match[1],
		}
	}

	stepsMatch := "@step_(\\w+)_(\\w+)_(\\w+)"

	match = regexp.MustCompile(stepsMatch).FindStringSubmatch(s0)
	if len(match) == 4 {

		if !isGoodLocation(match[2]) {
			return nil, false, fmt.Errorf("Location %s should be input or output", match[2])
		}

		source = &paramSource{
			Type:  sourceEvent,
			Name:  match[1],
			Loc:   match[2],
			Param: match[3],
		}
	}

	if source == nil {
		return nil, false, nil
	}

	return source, true, nil
}

func parseExternalSource(source *paramSource, variables *structs.JSONGeneric, jobs map[string]structs.Job) (string, bool, error) {
	// two types. @variable.ss or @steps.a.b.c.d
	if source.Type == sourceVariable {
		// variable
		inlinedValue, ok := (*variables)[source.Name]
		if !ok {
			return "", false, fmt.Errorf("The variable %s does not exists", source.Name)
		}

		// is an interface because is like parseInternalSource but in a variable
		return parseInternalSource(inlinedValue)
	} else if source.Type == sourceInput {

		// no se tiene que considerar, por eso el false
		return "", false, nil
	} else if source.Type == sourceEvent {
		// get the data from another step

		// get the job
		job, ok := jobs[source.Name]
		if !ok {
			return "", false, fmt.Errorf("Job specified %s does not exists", source.Name)
		}

		// get all the events given the location (input or output)
		var events structs.JSON
		switch source.Loc {
		case "input":
			events = job.Input
		case "output":
			events = job.Output
		default:
			panic("should not happen, found name " + source.Loc)
		}

		event, ok := events[source.Param]
		if !ok {
			return "", false, fmt.Errorf("Event name %s do not exists on job %s", source.Param, source.Name)
		}

		return event, true, nil
	}

	return "", false, fmt.Errorf("Type %d not found", source.Type)
}

func parseInternalSource(s interface{}) (string, bool, error) {
	// source of the form: {} or [] which are defined directly on the variables
	res, err := types.DetectType(s)
	return res, true, err
}

func inlineStep(step *structs.Step, variables *structs.JSONGeneric, jobs map[string]structs.Job) (map[string]string, map[string]bool, error) {
	params := map[string]string{}
	inputParams := map[string]bool{}

	// recorrer cada parametro del step
	for name, param := range step.Param {
		externalSource, ok, errr := parseParamSource(param)
		if errr != nil {
			return params, inputParams, fmt.Errorf("ERror parsing paramsoruce, %s", errr.Error())
		}

		var res string
		var use bool
		var err error

		if ok {
			// gestionar el external source
			res, use, err = parseExternalSource(externalSource, variables, jobs)
		} else {
			// gestionar el internal source
			res, use, err = parseInternalSource(param)
		}

		if err != nil {
			return params, inputParams, fmt.Errorf("Error parsing external/internal sources %s", err.Error())
		}

		if use {
			params[name] = res
		} else {
			inputParams[name] = true
		}

	}

	return params, inputParams, nil
}

func ValidateWorkflow(workflow *structs.Workflow, jobs *[]structs.Job) error {
	mapjobs := jobsToMap(*jobs)

	// Comprobar cada step
	for _, i := range workflow.Steps {

		// Sacar trabajo al que referencia el step
		job, ok := mapjobs[i.Job]
		if !ok {
			return fmt.Errorf("El job que requiere el step no esta en la lista de jobs. Error interno. ")
		}

		// inline del step para normalizar los valores
		params, inputParams, err := inlineStep(&i, &workflow.Variables, mapjobs)
		if err != nil {
			return err
		}

		// recorrer todos los inputs del job (obligatorios) y buscar el parametro
		for inputName, inputType := range job.Input {
			stepType, ok := params[inputName]
			if !ok {
				if _, ok := inputParams[inputName]; ok {
					// el parametro se especifica en la entrada.
					continue
				}
				return fmt.Errorf("El job %s necesita el input %s que no existe en el workflow", i.Job, inputName)
			}

			if stepType != inputType {
				return fmt.Errorf("Input do not match on %s. Necessary %s and found %s", i.Job, inputType, stepType)
			}
		}
	}

	return nil
}

func parseAllocExternalSource(source *paramSource, variables *structs.JSONGeneric, inputs *structs.JSONGeneric, previousAllocs *map[string]structs.Alloc) (interface{}, error) {
	// two types. @variable.ss or @steps.a.b.c.d
	if source.Type == sourceVariable || source.Type == sourceInput {

		fmt.Println("-- source type --")
		fmt.Println(source.Type)

		var inlinedValue structs.JSONGeneric
		if source.Type == sourceVariable {
			inlinedValue = *variables
		} else if source.Type == sourceInput {
			inlinedValue = *inputs
		}

		// variable
		res, ok := inlinedValue[source.Name]
		if !ok {
			return "", fmt.Errorf("The variable %s does not exists", source.Name)
		}

		// is an interface because is like parseInternalSource but in a variable
		return res, nil
	} else if source.Type == sourceEvent {

		// get the alloc
		alloc, ok := (*previousAllocs)[source.Name]
		if !ok {
			return "", fmt.Errorf("Job specified %s does not exists", source.Name)
		}

		// get all the events given the location (input or output)
		var events structs.JSONGeneric
		switch source.Loc {
		case "input":
			events = alloc.Input
		case "output":
			events = alloc.Output
		default:
			panic("should not happen, found name " + source.Loc)
		}

		event, ok := events[source.Param]
		if !ok {
			return "", fmt.Errorf("Event name %s do not exists on job %s", source.Param, source.Name)
		}

		return event, nil
	}

	return "", fmt.Errorf("Type %d not found", source.Type)
}

func CreateAllocFromStep(workflowAlloc *structs.WorkflowAlloc,
	step *structs.Step,
	jobs *map[string]structs.Job,
	variables *structs.JSONGeneric,
	userinputs *structs.JSONGeneric,
	previousAllocs *map[string]structs.Alloc) (*structs.Alloc, error) {

	job, ok := (*jobs)[step.Job]
	if !ok {
		return nil, fmt.Errorf("Job not found")
	}

	if previousAllocs == nil {
		// porque es la primera
		previousAllocs = &map[string]structs.Alloc{}
	}

	// los inputs de la ejecucion
	inputs := map[string]interface{}{}
	for inputname := range job.Input {
		param, ok := step.Param[inputname]
		if !ok {
			return nil, fmt.Errorf("No ha encontrado en step el input de nombre %s", inputname)
		}

		// compose the values
		externalSource, ok, errr := parseParamSource(param)
		if errr != nil {
			return nil, fmt.Errorf("Composealloc. error parsing paramsoruce, %s", errr.Error())
		}

		var res interface{}
		var err error

		if ok {
			// gestionar el external source
			res, err = parseAllocExternalSource(externalSource, variables, userinputs, previousAllocs)
		} else {
			// gestionar el internal source
			res = param
		}

		if err != nil {
			return nil, fmt.Errorf("Composealloc. Error parsing external/internal sources %s", err.Error())
		}

		inputs[inputname] = res
	}

	fmt.Println("//// \\\\ ")
	fmt.Println(workflowAlloc.WorkflowID)
	fmt.Println(workflowAlloc.ID)

	// create now the allocation struct
	alloc := structs.Alloc{
		WorkflowID:      workflowAlloc.WorkflowID,
		WorkflowAllocID: workflowAlloc.ID,
		JobID:           step.Job,
		Input:           inputs,
		Output:          structs.JSONGeneric{},
	}

	return &alloc, nil
}
