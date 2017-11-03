package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	crand "crypto/rand"
)

func (s *Steps) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Steps. Type assertion .([]byte) failed")
	}

	var i Steps

	err := json.Unmarshal(source, &i)
	if err != nil {
		return fmt.Errorf("Steps. Cast .(map[string]string) failed. Error %s", err.Error())
	}

	*s = i
	return nil
}

func (s Steps) Value() (driver.Value, error) {
	// parsear contenido
	j, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("Steps. Value cast failed %s", err.Error())
	}

	return string(j), nil
}

func (js *JSON) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Json. Type assertion .([]byte) failed")
	}

	var i map[string]string

	err := json.Unmarshal(source, &i)
	if err != nil {
		return errors.New("Json. Cast .(map[string]string) failed")
	}

	*js = i
	return nil
}

func (js JSON) Value() (driver.Value, error) {
	// vacio. inicializarlo
	if len(js) == 0 {
		js = map[string]string{}
	}

	// parsear contenido
	j, err := json.Marshal(js)
	if err != nil {
		return nil, errors.New("Json. Value cast failed")
	}

	return string(j), nil
}

func (jg *JSONGeneric) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("JSONGeneric. Type assertion .([]byte) failed")
	}

	var i map[string]interface{}

	err := json.Unmarshal(source, &i)
	if err != nil {
		return errors.New("JSONGeneric. Cast .(map[string]interface) failed")
	}

	*jg = i
	return nil
}

func (jg JSONGeneric) Value() (driver.Value, error) {
	// parsear contenido
	j, err := json.Marshal(jg)
	if err != nil {
		return nil, fmt.Errorf("JSONGeneric. Value cast failed %s", err.Error())
	}

	return string(j), nil
}

func (jg JSONGeneric) ToByte() ([]byte, error) {
	return json.Marshal(jg)
}

func (t *Timestamp) Scan(src interface{}) error {
	if src == nil {
		(*t) = 0
		return nil
	}

	source, ok := src.(int64)
	if !ok {
		return errors.New("Timestamp. Type assertion .([]byte) failed")
	}

	(*t) = Timestamp(source)
	return nil
}

func GenerateUUID() string {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}
