package logrus_support

import (
	"context"
	"testing"

	"github.com/ragpanda/go-toolkit/log"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, &LoggerTestSuite{})
}
func (suite *LoggerTestSuite) TestSuccess() {
	l := NewLogrusLogger(context.Background(), nil)
	Init(l, nil)

	ctx := context.Background()
	log.Debug(ctx, "i am first ")
	log.Info(ctx, "i am first ")
	log.Warn(ctx, "i am first ")
	log.Error(ctx, "i am first ")
	log.Fatal(ctx, "i am first ")
}
