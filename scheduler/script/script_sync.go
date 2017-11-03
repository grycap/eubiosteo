package script

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"pergamo/pergamo/structs"
	"pergamo/pergamo/types"
	"strings"
)

// Script sync shares with the volume the whole onedata space

type Sync struct {
}

type SyncState struct {
	Input, AllocID, DockerImage string
	InputMV, OutputMV           string
	ExportIni, ExportImage      string
	Zip                         string
	Port                        int
	Addrs                       string
}

var syncTemp = `#!/bin/bash

# Cambiar a running state
curl -XPOST http://{{.Addrs}}/allocs/{{.AllocID}}/running

# Compartir el fichero de input
mkdir -p /tmp/share
echo '{{.Input}}' > /tmp/share/input.json

# 1. Onedata -> /tmp/share (input images)
{{.InputMV}}

# 3. Export variables de los outputs
{{.ExportIni}}

# Ejecutar (Enlazar los export de los outputs...)
docker run --name {{.AllocID}} {{.ExportImage}} -v /tmp/share:/share {{.Zip}} {{.DockerImage}}

# 2. /tmp/share (output) -> Onedata
{{.OutputMV}}

# Logs a local
docker logs {{.AllocID}} > /tmp/stdout.log 2>/tmp/stderr.log

# Leer output y error (a variable)
OUPUT=$(cat /tmp/stdout.log)
ERROR=$(cat /tmp/stderr.log)

# Leer ficheros del output
OUTPUTJSON=$(cat /tmp/share/output.json)

# Enviar por curl los resultados
curl -XPOST http://{{.Addrs}}/allocs/{{.AllocID}}/done --data "logoutput=$OUPUT" --data "logerror=$ERROR" --data "jsonoutput=$OUTPUTJSON"

# Eliminar los tmp de resultados
# rm /tmp/stdout.log /tmp/stderr.log /tmp/share/output.json /tmp/share/input.json

# Eliminar contenedor
docker rm {{.AllocID}}
`

const shared = "/tmp/share"

func moveFiles(files structs.JSON, source, dest string, entry structs.JSONGeneric) (string, error) {
	moves := ""

	for inputName, inputType := range files {
		if types.IsFileType(inputType) {
			inputValue, ok := entry[inputName]
			if !ok {
				return "", errors.New("Fallo en move files porque el entry del job no esta en alloc. name => " + inputName)
			}

			inputStr, ok := inputValue.(string)
			if !ok {
				return "", errors.New("Failed to cast to string the input " + inputType)
			}

			moves += fmt.Sprintf("cp %s %s\n", filepath.Join(source, inputStr), filepath.Join(dest, inputStr))
		}
	}

	return moves, nil
}

func (c *Sync) Plan(onedata string, job *structs.Job, alloc *structs.Alloc, addrs string) (string, error) {
	input, err := createInputJson(job, alloc)
	if err != nil {
		return "", err
	}

	// Check zip sincronization
	Zip := ""
	if job.Zip != "" {
		Zip = "-v " + filepath.Join(onedata, job.Zip) + ":/home"
	}

	fmt.Println("- sincronizacion del zip -")
	fmt.Println(" sincronizando cosas ")

	// InputMV. input mv (1.)
	InputMV, err := moveFiles(job.Input, onedata, shared, alloc.Input)
	if err != nil {
		return "", err
	}

	// OutputMV. output mv (2.)
	OutputMV, err := moveFiles(job.Output, shared, onedata, alloc.Output)
	if err != nil {
		return "", err
	}

	// ExportIni. export output images (3.)
	ExportIni := ""

	// ExportImage. export output images dockerimage (4.)
	ExportImage := ""

	fmt.Println("-- output mv --")
	fmt.Println(OutputMV)

	fmt.Println("-- outputs --")
	fmt.Println(job.Output)
	fmt.Println(alloc.Output)

	for outputName, outputType := range job.Output {
		if types.IsFileType(outputType) {
			outputValue, ok := alloc.Output[outputName]
			if !ok {
				return "", errors.New("2. Fallo en el output porque el input del job no esta en alloc. name " + outputName)
			}

			ExportIni += fmt.Sprintf("export %s=%s\n", outputName, outputValue)
			ExportImage += fmt.Sprintf("%s=$%s", outputName, outputName)
		}
	}

	if ExportImage != "" {
		ExportImage = "-e " + ExportImage
	}

	jobdescription := SyncState{
		Input:       input,
		AllocID:     alloc.ID,
		DockerImage: job.DriverImage,
		InputMV:     InputMV,
		OutputMV:    OutputMV,
		ExportIni:   ExportIni,
		ExportImage: ExportImage,
		Zip:         Zip,
		Addrs:       addrs,
	}

	// Aplicar el template

	tmpl, err := template.New("jobdesc").Parse(syncTemp)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, jobdescription)
	if err != nil {
		return "", err
	}

	desc := buf.String()
	desc = strings.Replace(desc, "&#34;", "\"", -1) //	Porque en el template se annade eso cuando hay comillas

	return desc, nil
}

func NewSync() (Script, error) {
	s := Sync{}
	return &s, nil
}
