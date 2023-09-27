package injector

import (
	"context"
	"reflect"
	"sync"
)

type RegisterContext struct {
}

type RegisterProduct[T any] struct {
	Product T
	Cleaner func(ctx context.Context, rctx *RegisterContext)
}
type RegisterFunc[T any] func(ctx context.Context, rctx *RegisterContext) (RegisterProduct[T], error)

type Container struct {
	ctx context.Context

	register sync.Map // map[reflect.Type]RegisterFunc
	instance sync.Map // map[reflect.Type]interface{}

	instanceClean []func(ctx context.Context)
}

func NewContainer(ctx context.Context) (context.Context, *Container) {
	container := &Container{
		ctx:           ctx,
		register:      sync.Map{},
		instance:      sync.Map{},
		instanceClean: []func(ctx context.Context){},
	}

	return context.WithValue(ctx, defaultInjectContainerKey, container), container
}

//  Register for injection when needed
//  register: func(ctx context.Context) (t T, cleaner func(ctx context.Context)) {
//		return &T{}, nil
//	}
func (l *Container) Register(register RegisterFunc[any]) *Container {
	typeFlag := reflect.TypeOf(register).Out(0)
	l.register.Store(typeFlag, register)
	_, exist := l.instance.Load(typeFlag)
	if exist {
		l.instance.Delete(typeFlag)
	}

	return l
}

// Clear
func (l *Container) Clean(ctx context.Context) *Container {
	for i := len(l.instanceClean) - 1; i >= 0; i-- {
		if l.instanceClean[i] == nil {
			continue
		}
		l.instanceClean[i](ctx)
	}
	return l
}

// Example:
//		c.Register(func(ctx context.Context)(*A, func(ctx context.Context)){
//			return &A{}, nil
//		})
//
// 		a := &A{}
//		c.Inject(&a)
func (l *Container) Inject(fillPtr interface{}) error {
	instanceValue := reflect.ValueOf(fillPtr)
	if !instanceValue.IsValid() {
		panic("Instance is invalid")
	}

	instanceType := reflect.TypeOf(fillPtr).Elem()

	instance, exist := l.instance.Load(instanceType)
	if !exist {
		_instanceInit, exist := l.register.Load(instanceType)
		if !exist {
			return NewInjectorError("Type %s not register", instanceType.Name())
		}
		result := reflect.ValueOf(_instanceInit).Call([]reflect.Value{reflect.ValueOf(l.ctx)})
		newInstance, cleaner := result[0].Interface(), result[1].Interface().(func(context.Context))
		l.instance.Store(instanceType, newInstance)
		l.instanceClean = append(l.instanceClean, cleaner)
		instance = newInstance
	}
	instanceV := reflect.ValueOf(instance)
	instanceValue.Elem().Set(instanceV)

	return nil
}

func (l *Container) MustInject(fillPtr interface{}) {
	err := l.Inject(fillPtr)
	if err != nil {
		panic(err)
	}
}

func GetContainer(ctx context.Context) (*Container, error) {
	containerInf := ctx.Value(defaultInjectContainerKey)
	if containerInf == nil {
		return nil, NewInjectorError("Not find container")
	}

	container, ok := containerInf.(*Container)
	if !ok {
		return nil, NewInjectorError("Container type assert error, actually: %+v", containerInf)
	}

	return container, nil
}
