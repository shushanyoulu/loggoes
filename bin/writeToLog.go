//日志打印

package main

import (
	"log"
	"os"
)
var setPath = logPath()
func goeslog() *log.Logger {
	file, err := os.Create(setPath)
	if err != nil {
		log.Fatalln("fail to create goes.log file!")
	}
	logger := log.New(file, "[GOES] ", log.Lshortfile|log.LstdFlags)
	return logger
}
