package logs

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	line_info_inner_index = 8
	line_info_outer_index = 6
	log_file_ext          = ".log"
	log_file_fmt          = "20060102_1504"
	log_time_fmt          = "15:04:05.000"
	mem_log_salt          = "y!MIueW#@D"

	MEM_LOG_PFX = "^5#Mgqdxj&893o#hI3218QVf!GP13LKQ"
)

const (
	LEVEL_NONE = iota
	LEVEL_FATAL
	LEVEL_ERR
	LEVEL_WARNING
	LEVEL_INFO
	LEVEL_DEBUG
	LEVEL_STDOUT
)

var levelStrColorMap map[int]string = map[int]string{
	LEVEL_FATAL:   "\033[0;35mF \033[0m",
	LEVEL_ERR:     "\033[0;31mE \033[0m",
	LEVEL_WARNING: "\033[0;33mW \033[0m",
	LEVEL_INFO:    "I ",
	LEVEL_DEBUG:   "D ",
}

var levelStrMap = map[int]string{
	LEVEL_FATAL:   "F ",
	LEVEL_ERR:     "E ",
	LEVEL_WARNING: "W ",
	LEVEL_INFO:    "I ",
	LEVEL_DEBUG:   "D ",
}

const (
	LOG_KIND_STDOUT = iota
	LOG_KIND_FILE
	LOG_KIND_MEM
	LOG_KIND_CALL
	LOG_KIND_SIMPLE
)

var MemLogSaveFile string
var Timezone int = 8

var defaultLogger *Logger

type Option func(l *Logger)

func OptKind(kind int) Option {
	return func(l *Logger) {
		l.kind = kind
	}
}

func OptDelHour(hour int) Option {
	return func(l *Logger) {
		l.DelHour = hour
	}
}

func OptDelNum(num int) Option {
	return func(l *Logger) {
		l.DelNum = num
	}
}

func OptRotate(size int) Option {
	return func(l *Logger) {
		l.RotateSize = size
	}
}

func OptDir(dir string) Option {
	return func(l *Logger) {
		l.dir = dir
	}
}

func OptFile(name string) Option {
	return func(l *Logger) {
		l.fileName = name
	}
}

func OptWithColor(withColor bool) Option {
	return func(l *Logger) {
		l.WithColor = withColor
	}
}

func OptErrWithCode(errWithCode bool) Option {
	return func(l *Logger) {
		l.ErrWithCode = errWithCode
	}
}

func Start(levelStr string, options ...Option) {
	time.Local = time.FixedZone(fmt.Sprintf("UTC%+d", Timezone), Timezone*60*60)
	defaultLogger = New(levelStr, options...)
}

func New(levelStr string, options ...Option) *Logger {
	level := ParseLevel(levelStr)
	if level == LEVEL_NONE {
		return nil
	}

	logger := new(Logger)
	logger.level = level

	if logger.level == LEVEL_STDOUT {
		logger.kind = LOG_KIND_STDOUT
	} else {
		logger.kind = LOG_KIND_FILE
	}

	log.SetFlags(log.Ldate | log.Lmicroseconds)

	logger.dir = "."
	logger.fileName = "log"
	logger.DelHour = 72
	logger.DelNum = 10
	logger.RotateSize = 10
	logger.ErrWithCode = true
	logger.WithColor = true

	for _, v := range options {
		v(logger)
	}

	if logger.kind == LOG_KIND_FILE {
		logger.createLogFile()
	}

	return logger
}

func SetDefault(logger *Logger) {
	defaultLogger = logger
}

func GetDefault() *Logger {
	return defaultLogger
}

func Tag(tag string) *Logger {
	if defaultLogger == nil {
		return nil
	}

	defaultLogger.tag = tag
	return defaultLogger
}

func Debug(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Debug(v...)
	}
}

func Debugf(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Debugf(format, v...)
	}
}

func Debugj(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Debugj(v...)
	}
}

func Debugji(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Debugji(v...)
	}
}

func Info(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Info(v...)
	}
}

func Infof(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Infof(format, v...)
	}
}

func Infoj(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Infoj(v...)
	}
}

func Infoji(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Infoji(v...)
	}
}

func Warning(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Warning(v...)
	}
}

func Warningf(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Warningf(format, v...)
	}
}

func Warningj(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Warningj(v...)
	}
}

func Warningji(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Warningji(v...)
	}
}

func Error(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Error(v...)
	}
}

func Errorf(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Errorf(format, v...)
	}
}

func Errorj(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Errorj(v...)
	}
}

func Errorji(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Errorji(v...)
	}
}

func Fatal(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Fatal(v...)
	}
}

func Fatalf(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Fatalf(format, v...)
	}
}

func Fatalj(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Fatalj(v...)
	}
}

func Fatalji(v ...any) {
	if defaultLogger != nil {
		defaultLogger.Fatalji(v...)
	}
}

func PrintStack() {
	buf := debug.Stack()
	Error(string(buf))
}

func Line() string {
	return lineInfo(line_info_outer_index)
}

func DecryptLog(data []byte) {
	DecryptLogWithSalt(data, mem_log_salt)
}

func DecryptLogWithSalt(data []byte, salt string) {
	saltIndex := 0
	for k, v := range data {
		data[k] = v ^ salt[saltIndex]
		saltIndex = (saltIndex + 1) % len(salt)
	}
}

func ParseLevel(level string) int {
	switch level {
	case "stdout":
		return LEVEL_STDOUT
	case "debug":
		return LEVEL_DEBUG
	case "info":
		return LEVEL_INFO
	case "warn":
		return LEVEL_WARNING
	case "err":
		return LEVEL_ERR
	case "fatal":
		return LEVEL_FATAL
	}

	return LEVEL_NONE
}

func lineInfo(index int) string {
	buf := make([]byte, 1024)
	length := runtime.Stack(buf, false)
	arr := strings.Split(string(buf[:length]), "\n")
	if len(arr) > index {
		return arr[index]
	}

	return ""
}
