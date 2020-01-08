package common

import (
	"fmt"

	"github.com/alecthomas/log4go"
)

var Loger log4go.Logger

func init() {
	fmt.Print("log init\n")

	cfg := GetConfig()

	Loger = make(log4go.Logger)

	levelMap := map[string]log4go.Level{"DEBUG": log4go.DEBUG, "TRACE": log4go.TRACE, "INFO": log4go.INFO, "WARN": log4go.WARNING, "ERROR": log4go.ERROR, "CRITICAL": log4go.CRITICAL}
	level, _ := levelMap[cfg.Log.Log4Go.Level]
	logFileName := cfg.Log.FileName
	if len(logFileName) == 0 {
		logFileName = "./log/lottery.log"

	}
	if cfg.Log.Log4Go.Output == "stdout" {
		Loger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())

	}
	fileWriter := log4go.NewFileLogWriter(logFileName, true)
	fileWriter.SetRotate(true)
	fileWriter.SetRotateSize(cfg.Log.Log4Go.MaxSize)
	fileWriter.SetRotateDaily(false)
	fileWriter.SetRotateMaxBackup(20)
	Loger.AddFilter("file", level, fileWriter)
	Loger.Info("init log")

}
