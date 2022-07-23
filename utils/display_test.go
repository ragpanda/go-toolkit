package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DisplayTestSuite struct {
	suite.Suite
}

func TestDisplaySuite(t *testing.T) {
	suite.Run(t, &DisplayTestSuite{})
}
func (suite *DisplayTestSuite) TestSuccess() {
	ctx := context.Background()
	suite.T().Logf("%s", Display(ctx))
	suite.T().Logf("%s", Display(1))
	suite.T().Logf("%s", Display(&DisplayTestSuite{}))

	type Tmp struct {
		Name string
	}
	suite.T().Logf("%s", Display(&Tmp{}))
	suite.T().Logf("%s", MixUpDisplay(&Tmp{}, 1))
	suite.T().Logf("%s", MixUpDisplay(&Tmp{}, 0.5))
	suite.T().Logf("%s", MixUpDisplay(&Tmp{}, 0))
	suite.T().Logf("%s", DigestDisplay(&Tmp{}))

	t := &Tmp{}
	suite.Equal(MixUpDisplay(t, 0), Display(t))
	suite.Equal(DigestDisplay(t), DigestDisplay(t))
}
