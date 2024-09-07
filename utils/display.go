package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/cast"
)

// display object info, the struct will be convert to json, %+v as the final method, basic type will be use %+v too
func Display(data interface{}) string {
	dataValue := reflect.Indirect(reflect.ValueOf(data))
	if !dataValue.IsValid() {
		return fmt.Sprintf("%+v", data)
	}

	kind := dataValue.Type().Kind()
	if kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Map {
		if dataValue.Type() == reflect.TypeOf(sync.Map{}) {
			syncMap := data.(sync.Map)
			showStr := []string{}
			syncMap.Range(func(key, value interface{}) bool {
				showStr = append(showStr, fmt.Sprintf("\"%v\":%s", key, Display(value)))
				return true
			})
			return fmt.Sprintf("{%s}", strings.Join(showStr, ","))
		}

		b, _ := json.Marshal(data)
		result := string(b)
		if result == "" {
			result = fmt.Sprintf("%+v", data)
		}
		return result

	} else {
		result, err := cast.ToStringE(dataValue.Interface())
		if err != nil {
			return fmt.Sprintf("%+v", dataValue.Interface())
		}
		return result
	}

}

// mix up sensitive info to display
func MixUpDisplay(data interface{}, probability float32) string {
	dataStr := Display(data)
	dataBytes := []byte(dataStr)
	for i := range dataBytes {
		if rand.Float32() >= probability {
			dataBytes[i] = byte('*')
		}
	}
	return string(dataBytes)
}

// 简单摘要, 用于debug对比敏感信息
func DigestDisplay(data interface{}) string {
	plainText := Display(data)
	md5Obj := sha1.New()
	_, err := md5Obj.Write([]byte(plainText))
	if err != nil {
		return fmt.Sprintf("[Digest Failed size=%d]", len(plainText))
	}
	digest := md5Obj.Sum(nil)
	return hex.EncodeToString(digest)
}
