package common

import (
	mylog "mylog2"
)

var Logger *mylog.SimpleLogger

func InitLog(prefix string) {

	Logger = mylog.NewSimpleLogger(prefix)

}
