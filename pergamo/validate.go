package pergamo

import (
	"fmt"
	"os"
	"path/filepath"
	"pergamo/pergamo/structs"
	"pergamo/pergamo/types"
)

// AllocResponse comprueba que el output del alloc sea bueno
func (s *Server) AllocResponse(job *structs.Job, alloc *structs.Alloc) error {
	return nil
}

func ValidateOutput(onedata string, job *structs.Job, alloc *structs.Alloc) error {
	fmt.Println("-- Validate output --")
	fmt.Println(job)

	return validateAlloc(onedata, job.Output, &alloc.Output, job.Checkfiles)
}

func ValidateInput(onedata string, job *structs.Job, alloc *structs.Alloc) error {
	fmt.Println("-- Validate input --")
	fmt.Println(job)

	return validateAlloc(onedata, job.Input, &alloc.Input, job.Checkfiles)
}

// ValidateAlloc comprubea que los inputs sean buenos
func validateAlloc(onedata string, jobEvents structs.JSON, alloc *structs.JSONGeneric, checkfiles bool) error {

	for name, jobEvent := range jobEvents {
		inputEvent, ok := (*alloc)[name]
		if !ok {
			fmt.Println(jobEvents)
			fmt.Println(alloc)
			return fmt.Errorf("0. ValidateAlloc. La entrada de nombre %s no esta", name)
		}

		jobType, err := types.NewType(jobEvent)
		if err != nil {
			return fmt.Errorf("01. ValidateAlloc. La entrada de nombre %s no tiene el bien de codigo %s", name, jobEvent)
		}

		fmt.Println("- entrando -")

		if types.IsFileType(jobEvent) {

			// input/output es variable de file
			fmt.Println("EVENTO " + jobEvent + " es de los de file")
			fmt.Println(inputEvent)

			if !checkfiles {
				fmt.Println("No se validan los archivos. pasando")
				continue
			}

			inputEventStr, ok := inputEvent.(string)
			if !ok {
				return fmt.Errorf("1. ValidateAlloc. Input event for filetype cannot be casted to string")
			}

			fmt.Println("forma string")
			fmt.Println(inputEventStr)

			// inputEventStr should be the name of the file
			fileAbsPath := filepath.Join(onedata, inputEventStr)

			fmt.Println("archivo total")
			fmt.Println(fileAbsPath)
			fmt.Println("jobtype")
			fmt.Println(jobType)

			if _, err := os.Stat(fileAbsPath); os.IsNotExist(err) {
				return fmt.Errorf("File %s with path %s does not exists", name, fileAbsPath)
			}

			err = jobType.Parse(fileAbsPath)

			fmt.Println("parse del file")
			fmt.Println(err)

			if err != nil {
				fmt.Println("-- error en validate alloc --")
				return fmt.Errorf("2. ValidateAlloc. Error al validar entry file con path %s y error %s", fileAbsPath, err.Error())
			}

		} else {
			fmt.Println("- otra -")
			// input/output es variable inlined
			err = jobType.Parse(inputEvent)
			if err != nil {
				return fmt.Errorf("3. ValidateAlloc. Error al validar input con parametro para %s con error %s", name, err.Error())
			}
		}

	}

	return nil
}
