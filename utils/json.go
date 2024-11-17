package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func MustJsonDecodeString(v string, out interface{}) {
	if err := Unmarshal([]byte(v), out); err != nil {
		panic(fmt.Sprintf("json decode error: err=%s, data=%s", err.Error(), v))
	}
}

func MustJsonEncodeString(v interface{}) string {
	d, err := Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("json encode error: err=%s, data=%+v", err.Error(), d))
	}
	return string(d)
}

func MustJsonEncodeBytes(v interface{}) []byte {
	d, err := Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("json encode error: err=%s, data=%+v", err.Error(), d))
	}
	return d
}

func JsonDeepCopy[T any](v T) *T {
	d := new(T)
	MustJsonDecodeString(MustJsonEncodeString(v), &d)
	return d
}
