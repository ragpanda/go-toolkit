package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockService struct {
	Name string
}

type MockInterface interface {
	GetName() string
}

func (m *MockService) GetName() string {
	return m.Name
}
func TestBasicRegisterAndInject(t *testing.T) {
	loader := NewLoader()

	t.Run("type inject", func(t *testing.T) {
		var mockSvc *MockService
		assert.Nil(t, loader.Register(func() *MockService {
			mockSvc = &MockService{
				Name: "good",
			}
			return mockSvc
		}))

		injectObj := &MockService{}
		assert.Nil(t, loader.InjectByType(&injectObj))
		assert.Equal(t, "good", injectObj.Name, "name should match")
		assert.Equal(t, mockSvc, injectObj, "should be the same instance")

	})

	t.Run("interface inject", func(t *testing.T) {
		var mockSvc MockInterface
		assert.Nil(t, loader.Register(func() MockInterface {
			mockSvc = &MockService{
				Name: "good",
			}
			return mockSvc
		}))

		var injectObj MockInterface
		assert.Nil(t, loader.InjectByType(&injectObj))
		assert.Equal(t, "good", injectObj.GetName(), "name should match")
		assert.Equal(t, mockSvc, injectObj, "should be the same instance")
	})
}
func TestInjectByFunc(t *testing.T) {
	loader := NewLoader()
	t.Run("func args inject", func(t *testing.T) {
		var mockSvc *MockService
		assert.Nil(t, loader.Register(func() *MockService {
			mockSvc = &MockService{
				Name: "good",
			}
			return mockSvc
		}))

		var injectObj *MockService
		assert.Nil(t, loader.InjectByFuncArgs(func(m *MockService) {
			injectObj = m
		}))
		assert.Equal(t, "good", injectObj.Name, "name should match")
		assert.Equal(t, mockSvc, injectObj, "should be the same instance")
	})
}
