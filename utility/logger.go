package utility

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	debugMode bool = true
	Logger *LoggerStruct = defaultLogger()
)

type LoggerStruct struct {
	debugMode bool
	logger *log.Logger
	waitGroup sync.WaitGroup
}

func (l *LoggerStruct) Debug(message... interface{}) {
	if l.debugMode {
		l.waitGroup.Add(1)
		go func() {
			l.log("DEBUG", message...)
			l.waitGroup.Done()
		}()
	}
}

func (l *LoggerStruct) Info(message... interface{}) {
	l.waitGroup.Add(1)
	go func() {
		l.log("INFO", message...)
		l.waitGroup.Done()
	}()
}

func (l *LoggerStruct) Warning(message... interface{}) {
	l.waitGroup.Add(1)
	go func() {
		l.log("WARNING", message...)
		l.waitGroup.Done()
	}()
}

func (l *LoggerStruct) Error(message... interface{}) {
	l.waitGroup.Add(1)
	go func() {
		l.log("ERROR", message...)
		l.waitGroup.Done()
		
	}()
	l.Wait()
	os.Exit(1)
}

func (l *LoggerStruct) Wait() {
	l.waitGroup.Wait()
}

func (l *LoggerStruct) SetDebugMode(debugMode bool) {
	l.debugMode = debugMode
}

func (l *LoggerStruct) log(mode string, message... interface{}) {
	var builder strings.Builder
	builder.WriteString(mode)
	builder.WriteString(": ")
	for i := 0; i < len(message); i++ {
		builder.WriteString(fmt.Sprint(message[i]))
	}
	l.logger.Println(builder.String())
}

func NewLogger(debugMode bool, logger *log.Logger) *LoggerStruct {
	return &LoggerStruct{debugMode: debugMode, logger: logger}
}

func defaultLogger() *LoggerStruct {
	logDirectory := CurrentWorkingDirectory.Join("log")
	if !logDirectory.IsExist() {
		err := os.Mkdir(string(logDirectory), 0755)
		if err != nil {
			log.Fatalln("Error: Fail to create log directory, program abort.")
		}
	}
	currentTime := time.Now()
	logFilename := fmt.Sprintf("Log_%d%0.2d%0.2d%0.2d%0.2d%0.2d.log", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	logFilePath := logDirectory.Join(logFilename)
	logFile, err := os.OpenFile(string(logFilePath), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln("Error: Fail to open " + string(logFilePath) + ", program abort.")
	}
	return NewLogger(debugMode, log.New(logFile, "", log.Ldate|log.Ltime))
}
