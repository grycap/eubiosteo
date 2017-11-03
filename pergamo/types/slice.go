package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// String to convert the value to string and be stored in db easily
func intSliceToInterface(i []int) []interface{} {
	b := []interface{}{}
	for _, k := range i {
		b = append(b, k)
	}
	return b
}

func stringSliceToInterface(i []string) []interface{} {
	b := []interface{}{}
	for _, k := range i {
		b = append(b, k)
	}
	return b
}

func sliceString(value string) ([]interface{}, error) {
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = strings.TrimPrefix(value, "[")
		value = strings.TrimSuffix(value, "]")

		return stringSliceToInterface(strings.Split(value, ",")), nil
	}

	return []interface{}{}, errors.New("No es una lista")
}

type LSlice struct {
	Values []interface{}
	Field  LValueType
}

func NewNumberSlice() LValue {
	slice := new(LSlice)
	slice.Field = LTNumber

	return slice
}

func (ls *LSlice) Encode() (string, error) {
	data, err := json.Marshal(ls)
	return string(data), err
}

func (ls *LSlice) Decode(s string) error {

	var data LSlice
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return err
	}

	(*ls) = data
	return nil
}

// Value es el valor que se le pasa en el archivo json
func (ls *LSlice) Value() interface{} {
	return ls.Values
}

// Parse to convert the value again to his struct
func (ls *LSlice) Parse(s0 interface{}) error {

	v := reflect.ValueOf(s0)
	var values []interface{}

	// String structures can be parsed with json unmarshall
	// Case with data of the type '[1,2,3]'
	if v.Kind() == reflect.String {
		s1, ok := s0.(string)
		if !ok {
			return fmt.Errorf("Cannot parse to kind string")
		}

		if ls.Field == LTNumber {
			var res []int
			err := json.Unmarshal([]byte(s1), &res)
			if err != nil {
				return fmt.Errorf("String to LTNumber slice error " + err.Error())
			}

			values = intSliceToInterface(res)
		}
	} else {

		// To be used directly with the interface, we must use mapstructure
		if ls.Field == LTNumber {
			var res []int

			err := mapstructure.Decode(s0, &res)
			if err != nil {
				return fmt.Errorf("Mapstructure %s", err.Error())
			}

			values = intSliceToInterface(res)
		}
	}

	if values == nil {
		return fmt.Errorf("No values found")
	}

	(*ls).Values = values
	return nil
}

// Type to return his LValueType
func (ls *LSlice) Type() LValueType {
	return LTSlice
}
