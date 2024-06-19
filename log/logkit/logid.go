package logkit

import (
	"context"
	"fmt"

	"github.com/ragpanda/go-toolkit/log"
	util "github.com/ragpanda/go-toolkit/utils/idgen"
)

var LogIDInstance util.IDGenerator

func init() {
	_LogIDInstance, err := util.NewStandardSnowflake(util.GenerateMachineIDByMac(14))
	if err != nil {
		log.Error(context.Background(), "new logid instance failed %s", err.Error())
	}
	LogIDInstance = _LogIDInstance
}

func GetLogID() string {
	return fmt.Sprintf("%d", LogIDInstance.GenerateID())
}
