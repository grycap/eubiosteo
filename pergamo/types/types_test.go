package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestWithFile(t *testing.T) {
	dat, err := ioutil.ReadFile("./fixtures/simple.json")
	if err != nil {
		t.Fatal(err)
	}

	var res map[string]interface{}
	err = json.Unmarshal(dat, &res)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(dat))
	fmt.Println(res)

	hello, ok := res["hello"]
	if !ok {
		t.Fatal("Not found hello")
	}

	fmt.Println(hello)

	sl1 := LSlice{
		Field: LTNumber,
	}

	fmt.Println(hello)

	err = sl1.Parse(hello)
	if err != nil {
		t.Fatal("Error al parse " + err.Error())
	}

	fmt.Println(sl1)
	fmt.Println(sl1.Value())

	enc, err := sl1.Encode()
	if err != nil {
		t.Fatal("Error en el parse" + err.Error())
	}

	fmt.Println(enc)

	sl2 := LSlice{}
	err = sl2.Decode(enc)
	if err != nil {
		t.Fatal("DD" + err.Error())
	}

	fmt.Println(sl2.Values)
}

func TestNumberSlice(t *testing.T) {

	slice, err := NewType("slice.number")
	if err != nil {
		t.Fatal(err)
	}

	slice.Parse("[1,2,3]")

	fmt.Println(slice.Value())

}

func TestImage(t *testing.T) {
	imagePNG, err := NewType("image.png")
	if err != nil {
		t.Fatal(err)
	}

	err = imagePNG.Parse("./fixtures/lena.png")
	if err != nil {
		t.Fatal(err)
	}

	err = imagePNG.Parse("c4ac38bc-fd6e-8fbb-80b6-1dac171d042f")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDetectType(t *testing.T) {

	type Check struct {
		Name    string
		Content interface{}
	}

	checks := []Check{
		{"bool", true},
		{"bool", "true"},
		{"number", 1},
		{"slice.number", []int{1, 1, 1}},
		{"slice.number", "[1,2,3]"},
	}

	for index, check := range checks {
		typeExpect, err := DetectType(check.Content)
		if err != nil {
			t.Fatal(fmt.Errorf("Valor en indice %d de tipo %s devolvio error", index, check.Name))
		}

		if typeExpect != check.Name {
			t.Fatal(fmt.Errorf("Check failed on index %d and found %s and expected %s", index, typeExpect, check.Name))
		}
	}
}
