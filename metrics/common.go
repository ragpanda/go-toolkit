package metrics

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/go-metrics"
)

type Label struct {
	Name  string
	Value string
}

func EmitCounter(name string, value float32, labels ...Label) {
	key := append([]string{name}, labelsToKeys(labels)...)
	metrics.IncrCounter(key, value)
}

func EmitTimer(name string, value time.Duration, labels ...Label) {
	key := append([]string{name}, labelsToKeys(labels)...)
	metrics.AddSample(key, float32(value.Nanoseconds()))
}

func EmitKey(name string, value float32, labels ...Label) {
	key := append([]string{name}, labelsToKeys(labels)...)
	metrics.EmitKey(key, value)
}

func MapToLabel(m map[string]string) []Label {
	labels := make([]Label, 0, len(m))
	for k, v := range m {
		labels = append(labels, Label{Name: k, Value: v})
	}
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Name < labels[j].Name
	})
	return labels
}

func StructFieldToLabel(s interface{}) []Label {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	labels := make([]Label, 0, v.NumField())
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanInterface() {
			labels = append(labels, Label{
				Name:  t.Field(i).Name,
				Value: toString(field.Interface()),
			})
		}
	}

	return labels
}

func labelsToKeys(labels []Label) []string {
	keys := make([]string, len(labels))
	for i, label := range labels {
		keys[i] = label.Value
	}
	return keys
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
