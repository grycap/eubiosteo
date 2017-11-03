package types

import (
	"fmt"
	"reflect"
)

func stringToBool(s string) (bool, bool) {
	if s == "true" {
		return true, true
	} else if s == "false" {
		return false, true
	}

	return false, false
}

type LBool bool

func NewBool() LValue {
	return new(LBool)
}

func (lb *LBool) Encode() (string, error) {
	if *lb {
		return "true", nil
	}
	return "false", nil
}

func (lb *LBool) Decode(s string) error {
	value, ok := stringToBool(s)
	if !ok {
		return fmt.Errorf("Cannot decode bool with value %s", s)
	}

	(*lb) = LBool(value)
	return nil
}

func (lb *LBool) Value() interface{} {
	return lb
}

// Parse to convert the value again to his struct
func (lb *LBool) Parse(s interface{}) error {
	v := reflect.ValueOf(s)

	switch v.Kind() {
	case reflect.String:
		// buscar true | false
		s0, ok := s.(string)
		if !ok {
			return fmt.Errorf("Cannot convert a known string interface to string %s", s)
		}

		// Use decode to detect the strings true or false
		err := lb.Decode(s0)
		if err != nil {
			return fmt.Errorf("Cannot decode the string as bool %s", err)
		}

	case reflect.Bool:
		s0, ok := s.(bool)
		if !ok {
			return fmt.Errorf("Cannot convert a known bool interface to bool %s", s)
		}

		(*lb) = LBool(s0)

	default:
		return fmt.Errorf("No es un tipo valido")
	}

	return nil
}

// Type to return his LValueType
func (lb *LBool) Type() LValueType {
	return LTBool
}
