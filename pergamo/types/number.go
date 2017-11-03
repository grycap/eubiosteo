package types

import (
	"fmt"
	"strconv"
)

type LNumber int

func NewNumber() LValue {
	return new(LNumber)
}

func (ln *LNumber) Encode() (string, error) {
	str := strconv.Itoa(int(*ln))
	return str, nil
}

func (ln *LNumber) Decode(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	(*ln) = LNumber(num)
	return nil
}

func (ln *LNumber) Value() interface{} {
	return ln
}

// Parse to convert the value again to his struct
func (ln *LNumber) Parse(s interface{}) error {
	v0, ok := s.(int)
	if ok {
		(*ln) = LNumber(v0)
		return nil
	}

	v1, ok := s.(int32)
	if ok {
		(*ln) = LNumber(v1)
		return nil
	}

	v2, ok := s.(int64)
	if ok {
		(*ln) = LNumber(v2)
		return nil
	}

	v3, ok := s.(float32)
	if ok {
		(*ln) = LNumber(int(v3))
		return nil
	}

	v4, ok := s.(float64)
	if ok {
		(*ln) = LNumber(int(v4))
		return nil
	}

	return fmt.Errorf("Cannot parse %s", s)
}

// Type to return his LValueType
func (ln *LNumber) Type() LValueType {
	return LTNumber
}
