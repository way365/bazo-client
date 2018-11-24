package util

import (
	"io"
	"log"
	"os"
)

func InitLogger() *log.Logger {

	performanceLogFile, err := os.OpenFile("PerformanceLoggerClient.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	wrt := io.MultiWriter(os.Stdout, performanceLogFile)
	log.SetOutput(wrt)
	return log.New(wrt, "INFO: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}
