package injector

import "context"

func Inject(ctx context.Context, fillPtr interface{}) error {
	container, err := GetContainer(ctx)
	if err != nil {
		return err
	}

	return container.Inject(fillPtr)
}

func MustInject(ctx context.Context, fillPtr interface{}) {
	err := Inject(ctx, fillPtr)
	if err != nil {
		panic(err)
	}
}
