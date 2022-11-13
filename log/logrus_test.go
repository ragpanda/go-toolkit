package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, &LoggerTestSuite{})
}
func (suite *LoggerTestSuite) TestSuccess() {
	ctx := context.Background()
	Debug(ctx, "i am first ")
	Info(ctx, "i am first ")
	Warn(ctx, "i am first ")
	Error(ctx, "i am first ")
	Fatal(ctx, "i am first ")
}
