package utils

import (
	"bytes"
	"encoding/json"
)

func Unmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	return dec.Decode(v)
}

func Marshal(v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	ret := buffer.Bytes()
	// golang's encoder would always append a '\n', so we should drop it
	if len(ret) > 0 && ret[len(ret)-1] == '\n' {
		ret = ret[:len(ret)-1]
	}
	return ret, nil
}
