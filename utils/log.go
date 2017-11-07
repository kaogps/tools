package utils

import (
	"github.com/kdada/tinygo/log"
)

// Logger 控制台日志器
var Logger log.Logger

func init() {
	var err error
	Logger, err = log.NewLogger("console", "")
	if err != nil {
		panic(err)
	}
}
