/*
@Group : Jerryhax
@Author : Jerryhax
@DateTime : 1/11/20 11:36 AM
@File : log_maintain
@Version : v0.0.1.0
*/

package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type Level int

type Logger struct {
	//lock. when maintaining, cannot write
	sync.RWMutex
	//extend logger
	logger *log.Logger
	//log file.
	file     *os.File
	filePath string
	fileName string
}

var (
	L Logger

	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	LogRootPath        = "runtime/"
	TimeFormat         = "2006-01-02"
	LogFileExt         = "log"
	//每日0点0分重命名前一天的日志文件
	CronEveryday = "0 0 0 * * ?"
	//每月一日0点30分执行压缩打包
	CronEveryMonth = "0 30 0 1 * ?"
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup(filePath, fileName string) {
	var err error
	L.filePath = fmt.Sprintf("%s%s", LogRootPath, filePath)
	L.fileName = fmt.Sprintf("%s.%s", fileName, LogFileExt)

	L.file, err = MustOpen(L.fileName, L.filePath)
	if err != nil {
		log.Fatalln(err)
	}

	mw := io.MultiWriter(os.Stdout, L.file)
	L.logger = log.New(mw, DefaultPrefix, log.LstdFlags)
	Maintain()
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	L.logger.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	setPrefix(DEBUG)
	L.logger.Printf(format, v...)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	L.logger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	setPrefix(INFO)
	L.logger.Printf(format, v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	L.logger.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	setPrefix(WARNING)
	L.logger.Printf(format, v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	L.logger.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	setPrefix(ERROR)
	L.logger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	L.logger.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	setPrefix(FATAL)
	L.logger.Printf(format, v...)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], file, line)
		//logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	L.logger.SetPrefix(logPrefix)
}
