package types

import (
	"fmt"
	"strings"
)

type LValueType int

const (
	LTNumber LValueType = iota + 1
	LTSlice
	LTImage
	LTBool
)

var types = map[string]LValue{
	"bool":         NewBool(),
	"number":       NewNumber(),
	"slice.number": NewNumberSlice(),
	"image.jpeg":   NewJPEGImage(),
	"image.png":    NewPNGImage(),
	"image.nifti":  NewNIFTIImage(),
	"file.other":   NewOtherFile(),
}

func IsFileType(typename string) bool {
	if strings.HasPrefix(typename, "image") {
		return true
	}

	if strings.HasPrefix(typename, "file") {
		return true
	}

	return false
}

func GetFileExtension(typename string) (string, error) {
	if !IsFileType(typename) {
		return "", fmt.Errorf("is not file type")
	}

	typeDef, ok := types[typename]
	if !ok {
		return "", fmt.Errorf("not found on types")
	}

	image, ok := typeDef.(*LImage)
	if !ok {
		fmt.Println("----")
		fmt.Println(typeDef)
		return "", fmt.Errorf("cannot convert to LTImage")
	}

	extension, ok := imageExtensions[image.Format]
	if !ok {
		return "", fmt.Errorf("extension not found")
	}

	return extension, nil
}

// DetectType te da el codigo (como el de newtype sobre que tipo de valor es)
func DetectType(s interface{}) (string, error) {
	for name, lvalue := range types {
		err := lvalue.Parse(s)
		if err == nil {
			return name, nil
		}
	}

	return "", fmt.Errorf("No se ha encontrado nada para %s", s)
}

func NewType(s string) (LValue, error) {
	t, ok := types[s]
	if !ok {
		return nil, fmt.Errorf("No hay type para el value %s", s)
	}

	return t, nil
}

type LValue interface {
	// Encode para enviarlo a bd
	Encode() (string, error)

	// Decode para sacar el archivo desde bd
	Decode(string) error

	// Parse to convert the value again to his struct
	Parse(interface{}) error

	// Value es el valor que se pasa en los archivos json
	Value() interface{}

	// Type to return his LValueType
	Type() LValueType
}
