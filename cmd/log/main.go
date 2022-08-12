package main

import (
	"os"

	"github.com/pedrogao/plib/pkg/log"
)

func main() {
	log.Info("std log")
	log.Info("std log")
	log.SetOptions(log.WithLevel(log.DebugLevel))
	log.Debug("change std log to debug level")
	log.SetOptions(log.WithFormatter(&log.JsonFormatter{IgnoreBasicFields: false}))
	log.Debug("log in json format")
	log.Info("another log in json format") // 输出到文件

	fd, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("create file test.log failed")
	}

	defer fd.Close()

	l := log.New(log.WithLevel(log.InfoLevel), log.WithOutput(fd),
		log.WithFormatter(&log.JsonFormatter{IgnoreBasicFields: false}))
	l.Info("custom log with json formatter")
}
