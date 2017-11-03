package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type ImageFormat int

const (
	JPEG ImageFormat = iota + 1
	PNG
	NIFTI
	OTHER
)

func (i *ImageFormat) String() string {
	switch *i {
	case JPEG:
		return "JPEG"
	case PNG:
		return "PNG"
	case NIFTI:
		return "NIFTI"
	case OTHER:
		return "OTHER"
	}

	panic("Not found image format")
}

func NewJPEGImage() LValue {
	return NewImage(JPEG)
}

func NewPNGImage() LValue {
	return NewImage(PNG)
}

func NewNIFTIImage() LValue {
	return NewImage(NIFTI)
}

func NewOtherFile() LValue {
	return NewImage(OTHER)
}

func NewImage(format ImageFormat) LValue {
	image := new(LImage)
	image.Format = format

	return image
}

var imageExtensions = map[ImageFormat]string{
	JPEG:  "jpeg",
	PNG:   "png",
	NIFTI: "nii",
	OTHER: "",
}

var imageValidations = map[ImageFormat]func(content *[]byte) bool{
	JPEG: func(data *[]byte) bool {
		buf := *data
		return len(buf) > 2 &&
			buf[0] == 0xFF &&
			buf[1] == 0xD8 &&
			buf[2] == 0xFF
	},
	PNG: func(data *[]byte) bool {
		buf := *data
		return len(buf) > 3 &&
			buf[0] == 0x89 && buf[1] == 0x50 &&
			buf[2] == 0x4E && buf[3] == 0x47
	},
	NIFTI: func(data *[]byte) bool { //	Approx..
		buf := *data
		return len(buf) > 4 &&
			buf[0] == 0x5C && buf[1] == 0x01 &&
			buf[2] == 0x00 && buf[3] == 0x00
	},
	OTHER: func(data *[]byte) bool {
		return true
	},
}

func GetImageFormat(content *[]byte) string {
	for name, validate := range imageValidations {
		if validate(content) {
			return name.String()
		}
	}

	return "" // OTHER
}

func validateImage(format ImageFormat, content []byte) bool {
	validate, ok := imageValidations[format]
	if !ok {
		panic(fmt.Errorf("Pide una validacion de un formato de imagen que no existe %d", format))
	}

	return validate(&content)
}

type LImage struct {
	Path   string
	Format ImageFormat
}

func (li *LImage) Encode() (string, error) {
	data, err := json.Marshal(li)
	return string(data), err
}

func (li *LImage) Decode(s string) error {
	var lj LImage
	err := json.Unmarshal([]byte(s), &lj)
	if err != nil {
		return err
	}

	(*li) = lj
	return nil
}

func (li *LImage) Value() interface{} {
	return li.Path
}

func (li *LImage) Parse(s interface{}) error {
	fmt.Println("-- read image ## --")
	fmt.Println(s)

	path, ok := s.(string)
	if !ok {
		return fmt.Errorf("Debe ser string para hacer parse")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	valid := validateImage(li.Format, data)
	if !valid {
		return errors.New("No es valido el formato y la imagen. Se esperaba que fuera " + li.Format.String())
	}

	return nil
}

func (li *LImage) Type() LValueType {
	return LTImage
}
