package injector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ContainerTestSuite struct {
	suite.Suite
}

func TestContainerSuite(t *testing.T) {
	suite.Run(t, &ContainerTestSuite{})
}

type TestInjectType1 struct {
}

func (suite *ContainerTestSuite) TestInjectStructSuccess() {
	ctx := context.Background()
	ctx, c := NewContainer(ctx)
	t := &TestInjectType1{}
	c.Register(func(ctx context.Context) (*TestInjectType1, func(ctx context.Context)) {
		return t, nil
	})

	var v *TestInjectType1
	suite.Nil(Inject(ctx, &v))
	suite.Equal(t, v)
}

type TestInjectInf1 interface{}

func (suite *ContainerTestSuite) TestInjectInterfaceSuccess() {
	ctx := context.Background()
	ctx, c := NewContainer(ctx)
	t := &TestInjectType1{}

	// interface
	var tInf TestInjectInf1 = t
	c.Register(func(ctx context.Context) (TestInjectInf1, func(ctx context.Context)) {
		return tInf, nil
	})

	// struct ptr
	c.Register(func(ctx context.Context) (*TestInjectType1, func(ctx context.Context)) {
		return t, nil
	})

	var v TestInjectInf1
	err := Inject(ctx, &v)
	suite.Equal(tInf, v)

	suite.T().Log(err)
	suite.Nil(err)
}

func (suite *ContainerTestSuite) TestFail() {
	ctx := context.Background()
	ctx, _ = NewContainer(ctx)

	var v *TestInjectType1
	suite.NotNil(Inject(ctx, &v))
	suite.Nil(v)

}
