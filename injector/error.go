package injector

import "fmt"

type InjectError struct {
	err string
}

func NewInjectorError(errorStrFmt string, args ...interface{}) *InjectError {
	err := &InjectError{
		err: fmt.Sprintf(errorStrFmt, args...),
	}
	return err
}

func (self *InjectError) Error() string {
	return self.err
}
