package loader

import (
	"reflect"
	"sync"

	"go.uber.org/dig"
)

func NewLoader() *Loader {
	loader := &Loader{
		Container: dig.New(),
	}
	return loader
}

func InitGlobalLoader() *Loader {
	l := NewLoader()
	defaultLoader = l
	return l
}

func GetDefaultLoader() *Loader {
	initOnce.Do(func() {
		defaultLoader = NewLoader()
	})

	return defaultLoader
}

var initOnce sync.Once
var defaultLoader *Loader

type Loader struct {
	*dig.Container
}

func (self *Loader) Register(constructor interface{}, opts ...dig.ProvideOption) error {
	return self.Provide(constructor, opts...)
}

// InjectByType is support give a ptr of the type(include struct or interface)
// and loader will inject the instance to the ptr
func (c *Loader) InjectByType(object interface{}, opts ...dig.InvokeOption) error {
	f := reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{reflect.TypeOf(object).Elem()}, []reflect.Type{}, false),
		func(args []reflect.Value) (results []reflect.Value) {
			reflect.ValueOf(object).Elem().Set(args[0])
			return
		})
	return c.Invoke(f.Interface(), opts...)
}

// InjectByFuncArgs inject by func args
func (c *Loader) InjectByFuncArgs(function interface{}, opts ...dig.InvokeOption) error {
	return c.Invoke(function, opts...)
}
