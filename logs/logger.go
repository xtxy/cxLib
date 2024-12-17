package logger

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	line_info_inner_index = 8
	line_info_outer_index = 6
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
	LEVEL_NOTICE
	LEVEL_INFO
	LEVEL_DEBUG
	LEVEL_STDOUT
)

var levelString map[int]string = map[int]string{
	LEVEL_FATAL:   "\033[0;35mF \033[0m",
	LEVEL_ERR:     "\033[0;31mE \033[0m",
	LEVEL_WARNING: "\033[0;33mW \033[0m",
	LEVEL_NOTICE:  "N ",
	LEVEL_INFO:    "I ",
	LEVEL_DEBUG:   "D ",
}

const (
	log_kind_file = iota
	log_kind_stdout
	log_kind_mem
	log_kind_call
	log_kind_simple
)

type Logger struct {
	kind   int
	level  int
	logger *log.Logger

	// file kind
	dir         string
	filePfx     string
	file        *os.File
	timeout     int64
	currentSize int

	DelHour    int
	DelNum     int
	RotateSize int

	// mem kind
	buf      []string
	bufIndex int

	Salt string

	tag string

	// call kind
	CallFunc func(timeStr, data string)
}

var MemLogSaveFile string
var Timezone int = 8

var defaultLogger *Logger

type Option func(l *Logger)

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

func OptFilePfx(pfx string) Option {
	return func(l *Logger) {
		l.filePfx = pfx
	}
}

func Start(dir string, filePfx string, level int, options ...Option) {
	time.Local = time.FixedZone(fmt.Sprintf("UTC%+d", Timezone), Timezone*60*60)

	options = append(options, OptDir(dir), OptFilePfx(filePfx))
	defaultLogger = New(level, options...)
	defaultLogger.createLogFile()
}

func StartLevel(levelStr string, options ...Option) {
	level, create := ParseLevel(levelStr)
	if !create {
		return
	}

	time.Local = time.FixedZone(fmt.Sprintf("UTC%+d", Timezone), Timezone*60*60)
	defaultLogger = New(level, options...)
	defaultLogger.createLogFile()
}

func StartFile(name string, level int) {
	time.Local = time.FixedZone(fmt.Sprintf("UTC%+d", Timezone), Timezone*60*60)

	defaultLogger = new(Logger)
	defaultLogger.level = level
	defaultLogger.kind = log_kind_file
	defaultLogger.createLogFileByName(name)
}

func StartLogger(logger *Logger) {
	time.Local = time.FixedZone(fmt.Sprintf("UTC%+d", Timezone), Timezone*60*60)
	defaultLogger = logger
}

func GetDefaultLogger() *Logger {
	return defaultLogger
}

func NewStdout() *Logger {
	return newLogger(LEVEL_STDOUT, log_kind_stdout)
}

func NewConsole(level int) *Logger {
	return newLogger(level, log_kind_stdout)
}

func NewMem(level int) *Logger {
	return newLogger(level, log_kind_mem)
}

func NewCall(level int) *Logger {
	return newLogger(level, log_kind_call)
}

func newLogger(level int, kind int) *Logger {
	logger := new(Logger)
	logger.level = level
	logger.kind = kind
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	return logger
}

func New(level int, options ...Option) *Logger {
	logger := newLogger(level, log_kind_file)
	logger.dir = "."
	logger.filePfx = "log"
	logger.DelHour = 72
	logger.DelNum = 10
	logger.RotateSize = 10

	for _, v := range options {
		v(logger)
	}

	logger.createLoggerDir()

	return logger
}

func NewSimple(logger *log.Logger, level int) *Logger {
	myLogger := new(Logger)
	myLogger.logger = logger
	myLogger.level = level
	myLogger.kind = log_kind_simple

	return myLogger
}

func ParseLevel(level string) (int, bool) {
	number := parseLevel(level)
	if number != -1 {
		return number, true
	}

	if level == "stdout" {
		defaultLogger = NewStdout()
		number = LEVEL_STDOUT
	} else if strings.HasPrefix(level, "console") {
		number = parseLevel(strings.Split(level, "|")[1])
		defaultLogger = NewConsole(number)
	} else if strings.HasPrefix(level, "mem") {
		number = parseLevel(strings.Split(level, "|")[1])
		defaultLogger = NewMem(number)
		defaultLogger.RotateSize = 2000
		defaultLogger.Salt = mem_log_salt
	}

	return number, false
}

func parseLevel(level string) int {
	switch level {
	case "debug":
		return LEVEL_DEBUG
	case "info":
		return LEVEL_INFO
	case "notice":
		return LEVEL_NOTICE
	case "warn":
		return LEVEL_WARNING
	case "err":
		return LEVEL_ERR
	case "fatal":
		return LEVEL_FATAL
	}
	return -1
}

func PrintJson(level int, v ...any) {
	defaultLogger.PrintJson(level, v...)
}

func PrintJsonFmt(level int, v ...any) {
	defaultLogger.PrintJsonFmt(level, true, v...)
}

func Printf(level int, format string, v ...any) {
	if level == LEVEL_ERR {
		str := fmt.Sprintf(format, v...)
		defaultLogger.Error(str)
	} else {
		defaultLogger.Printf(level, format, v...)
	}
}

func Tag(tag string) *Logger {
	defaultLogger.tag = tag
	return defaultLogger
}

func Debug(v ...any) {
	defaultLogger.Debug(v...)
}

func Debugf(format string, v ...any) {
	Printf(LEVEL_DEBUG, format, v...)
}

func Debugj(v ...any) {
	PrintJson(LEVEL_DEBUG, v...)
}

func Info(v ...any) {
	defaultLogger.Info(v...)
}

func Infof(format string, v ...any) {
	Printf(LEVEL_INFO, format, v...)
}

func Infoj(v ...any) {
	PrintJson(LEVEL_INFO, v...)
}

func Notice(v ...any) {
	defaultLogger.Notice(v...)
}

func Noticef(format string, v ...any) {
	Printf(LEVEL_NOTICE, format, v...)
}

func Noticej(v ...any) {
	PrintJson(LEVEL_NOTICE, v...)
}

func Warning(v ...any) {
	defaultLogger.Warning(v...)
}

func Warningf(format string, v ...any) {
	Printf(LEVEL_WARNING, format, v...)
}

func Warningj(v ...any) {
	PrintJson(LEVEL_WARNING, v...)
}

func Error(v ...any) {
	defaultLogger.Error(v...)
}

func Errorf(format string, v ...any) {
	Printf(LEVEL_ERR, format, v...)
}

func Errorj(v ...any) {
	PrintJson(LEVEL_ERR, v...)
}

func Fatal(v ...any) {
	defaultLogger.Fatal(v)
}

func PrintStack() {
	buf := debug.Stack()
	Error(string(buf))
}

func Line() string {
	return lineInfo(line_info_outer_index)
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

func (l *Logger) GetStdLogger() *log.Logger {
	return l.logger
}

func (l *Logger) ChangeLevel(level int) {
	l.level = level
}

func (l *Logger) ChangeLevelByStr(levelStr string) bool {
	level := parseLevel(levelStr)
	if level <= 0 {
		return false
	}

	l.level = level
	return true
}

func (l *Logger) GetLevel() int {
	return l.level
}

func (l *Logger) createLogFileByName(name string) {
	if nil != l.file {
		l.file.Close()
		l.file = nil
	}

	os.Remove(name)

	var err error
	l.file, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("open log file failed, error:", err)
		return
	}

	l.timeout = time.Now().Add(time.Duration(l.DelHour) * time.Hour).Unix()
	if l.logger != nil {
		l.logger.SetOutput(l.file)
	} else {
		l.logger = log.New(l.file, "", log.Ldate|log.Lmicroseconds)
	}
}

func (l *Logger) createLogFile() {
	l.currentSize = 0

	logFileName := fmt.Sprintf("%s/%s_%s.log", l.dir, l.filePfx, time.Now().Format(log_file_fmt))
	l.createLogFileByName(logFileName)
	l.delOldLogs()
}

func (l *Logger) checkFileRotate() {
	if l.RotateSize == 0 && l.DelHour == 0 && l.DelNum == 0 || l.file == nil {
		return
	}

	if l.currentSize > l.RotateSize*1024*1024 || time.Now().Unix() >= l.timeout {
		l.createLogFile()
	}
}

func (l *Logger) delOldLogs() {
	if l.DelHour == 0 && l.DelNum == 0 {
		return
	}

	currentFileName := ""
	if nil != l.file {
		fileStat, err := l.file.Stat()
		if err != nil {
			Error("l.file.stat.error:", err)
		} else {
			currentFileName = fileStat.Name()
		}
	}

	type logFileInfo struct {
		name    string
		modTime int64
	}

	infos := make([]logFileInfo, 0)

	delTime := time.Now().Add(time.Duration(-l.DelHour) * time.Hour).Unix()
	filepath.Walk(l.dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			Error("logger.rmoldlog.walk.error:", err)
			return nil
		}

		if info.IsDir() || info.Name() == currentFileName {
			return nil
		}

		if !strings.HasPrefix(info.Name(), l.filePfx+"_") || !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}

		if info.ModTime().Unix() > delTime {
			infos = append(infos, logFileInfo{
				name: file, modTime: info.ModTime().Unix(),
			})

			return nil
		}

		err = os.Remove(file)
		return err
	})

	if l.DelNum > 0 && len(infos) > l.DelNum {
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].modTime > infos[j].modTime
		})

		for _, v := range infos[l.DelNum:] {
			os.Remove(v.name)
		}
	}
}

func (l *Logger) PrintJsonFmt(level int, needFmt bool, params ...any) {
	defer func() {
		if err := recover(); err != nil {
			l.Error(params...)
		}
	}()

	var slice []byte
	var err error
	newParams := make([]any, len(params))
	for k, v := range params {
		value := reflect.ValueOf(v)
		if !value.IsValid() || value.IsZero() {
			continue
		}

		switch reflect.ValueOf(v).Type().Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Ptr:
			if needFmt {
				slice, err = json.MarshalIndent(v, "", "    ")
			} else {
				slice, err = json.Marshal(v)
			}

			if err != nil {
				newParams[k] = v
			} else {
				newParams[k] = string(slice)
			}

		default:
			newParams[k] = v
		}
	}

	switch level {
	case LEVEL_DEBUG:
		l.Debug(newParams...)

	case LEVEL_INFO:
		l.Info(newParams...)

	case LEVEL_NOTICE:
		l.Notice(newParams...)

	case LEVEL_WARNING:
		l.Warning(newParams...)

	case LEVEL_ERR:
		l.Error(newParams...)

	case LEVEL_FATAL:
		l.Fatal(newParams...)
	}
}

func (l *Logger) Printf(level int, format string, v ...any) {
	str := fmt.Sprintf(format, v...)

	switch level {
	case LEVEL_DEBUG:
		l.Debug(str)

	case LEVEL_INFO:
		l.Info(str)

	case LEVEL_NOTICE:
		l.Notice(str)

	case LEVEL_WARNING:
		l.Warning(str)

	case LEVEL_ERR:
		l.Error(str)

	case LEVEL_FATAL:
		l.Fatal(str)
	}
}

func (l *Logger) PrintJson(level int, params ...any) {
	l.PrintJsonFmt(level, false, params...)
}

func (l *Logger) output(level int, v ...any) {
	if l.level < level {
		l.tag = ""
		return
	}

	tag := l.tag
	l.tag = ""

	str := fmt.Sprint(v)
	str = string([]byte(str)[1 : len(str)-1])
	levelStr := levelString[level]

	if tag != "" {
		str = levelStr + tag + ": " + str
	} else {
		str = levelStr + str
	}

	if len(str) > 4*1024*1024 {
		str = "too long:" + strconv.Itoa(len(str)) + ",content:" + string([]byte(str)[:100]) + "..."
	}

	switch l.kind {
	case log_kind_stdout:
		log.Println(str)

	case log_kind_file:
		if nil == l.file {
			l.createLogFile()
		}

		l.logger.Println(str)

		l.currentSize += len(log_time_fmt) + len(str) + 1
		l.checkFileRotate()

	case log_kind_simple:
		l.logger.Println(str)

	case log_kind_mem:
		if nil == l.buf {
			l.buf = make([]string, l.RotateSize)
			l.bufIndex = 0
		}

		l.buf[l.bufIndex] = "[" + time.Now().Format(log_time_fmt) + "]" + str
		l.bufIndex = (l.bufIndex + 1) % l.RotateSize

	case log_kind_call:
		if nil != l.CallFunc {
			l.CallFunc(time.Now().Format(log_time_fmt), str)
		}
	}
}

func (l *Logger) GetMemLog() []byte {
	if len(l.buf) == 0 {
		return nil
	}

	arr := make([]string, 0)
	for i := 0; i < l.RotateSize; i++ {
		index := (l.bufIndex + i + 1) % l.RotateSize
		str := l.buf[index]
		if str != "" {
			arr = append(l.buf, str)
		}
	}

	logStr := strings.Join(arr, "\n")

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write([]byte(logStr)); err != nil {
		return nil
	}
	if err := gz.Close(); err != nil {
		return nil
	}

	data := buf.Bytes()
	DecryptLog(data)

	return data
}

func (l *Logger) SaveMemLog() {
	data := l.GetMemLog()
	if len(data) == 0 {
		return
	}

	logStr := MEM_LOG_PFX + base64.StdEncoding.EncodeToString(data)

	if MemLogSaveFile == "" {
		MemLogSaveFile = "runLog.bin"
	}

	os.WriteFile(MemLogSaveFile, []byte(logStr), 0644)
}

func (l *Logger) Tag(tag string) *Logger {
	l.tag = tag
	return l
}

func (l *Logger) Debug(v ...any) {
	l.output(LEVEL_DEBUG, v...)
}

func (l *Logger) Info(v ...any) {
	l.output(LEVEL_INFO, v...)
}

func (l *Logger) Notice(v ...any) {
	l.output(LEVEL_NOTICE, v...)
}

func (l *Logger) Warning(v ...any) {
	l.output(LEVEL_WARNING, v...)
}

func (l *Logger) Error(v ...any) {
	line := lineInfo(line_info_inner_index)
	if line != "" {
		l.output(LEVEL_ERR, v...)
		l.output(LEVEL_ERR, line)
	} else {
		l.output(LEVEL_ERR, v...)
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	str := fmt.Sprintf(format, v...)
	l.Error(str)
}

func (l *Logger) Fatal(v ...interface{}) {
	line := lineInfo(line_info_inner_index)
	if line != "" {
		l.output(LEVEL_FATAL, v...)
		l.output(LEVEL_FATAL, line)
	} else {
		l.output(LEVEL_FATAL, v...)
	}

	panic(v)
}

func (l *Logger) createLoggerDir() {
	if l.dir == "." {
		return
	}

	dirInfo, err := os.Stat(l.dir)
	if nil == err && dirInfo.IsDir() {
		return
	}

	err = os.MkdirAll(l.dir, 0777)
	if err != nil {

		fmt.Println("createLoggerDir, err:", err)
		panic(err)
	}
}
