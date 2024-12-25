package logs

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	DelHour     int
	DelNum      int
	DelSize     int64
	RotateSize  int
	ErrWithCode bool
	WithColor   bool
	Salt        string

	// call kind
	CallFunc func(timeStr, data string)

	kind   int
	level  int
	logger *log.Logger

	// file kind
	dir         string
	fileName    string
	file        *os.File
	timeout     int64
	currentSize int

	// mem kind
	buf      []string
	bufIndex int

	tag string
}

func (l *Logger) GetGoLogger() *log.Logger {
	return l.logger
}

func (l *Logger) ChangeLevel(level int) {
	l.level = level
}

func (l *Logger) ChangeLevelByStr(levelStr string) bool {
	level := ParseLevel(levelStr)
	if level == LEVEL_NONE {
		return false
	}

	l.level = level
	return true
}

func (l *Logger) GetLevel() int {
	return l.level
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
	l.show(LEVEL_DEBUG, v)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.showf(LEVEL_DEBUG, format, v)
}

func (l *Logger) Debugj(v ...any) {
	l.showj(LEVEL_DEBUG, v)
}

func (l *Logger) Debugji(v ...any) {
	l.showji(LEVEL_DEBUG, v)
}

func (l *Logger) Info(v ...any) {
	l.show(LEVEL_INFO, v)
}

func (l *Logger) Infof(format string, v ...any) {
	l.showf(LEVEL_INFO, format, v)
}

func (l *Logger) Infoj(v ...any) {
	l.showj(LEVEL_INFO, v)
}

func (l *Logger) Infoji(v ...any) {
	l.showji(LEVEL_INFO, v)
}

func (l *Logger) Warning(v ...any) {
	l.show(LEVEL_WARNING, v)
}

func (l *Logger) Warningf(format string, v ...any) {
	l.showf(LEVEL_WARNING, format, v)
}

func (l *Logger) Warningj(v ...any) {
	l.showj(LEVEL_WARNING, v)
}

func (l *Logger) Warningji(v ...any) {
	l.showji(LEVEL_WARNING, v)
}

func (l *Logger) Error(v ...any) {
	l.show(LEVEL_ERR, v)

	if l.ErrWithCode {
		line := lineInfo(line_info_inner_index)
		if line != "" {
			l.output(l.combineOutputStr(LEVEL_ERR, line))
		}
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	l.showf(LEVEL_ERR, format, v)

	if l.ErrWithCode {
		line := lineInfo(line_info_inner_index)
		if line != "" {
			l.output(l.combineOutputStr(LEVEL_ERR, line))
		}
	}
}

func (l *Logger) Errorj(v ...any) {
	l.showj(LEVEL_ERR, v)

	if l.ErrWithCode {
		line := lineInfo(line_info_inner_index)
		if line != "" {
			l.output(l.combineOutputStr(LEVEL_ERR, line))
		}
	}
}

func (l *Logger) Errorji(v ...any) {
	l.showji(LEVEL_ERR, v)

	if l.ErrWithCode {
		line := lineInfo(line_info_inner_index)
		if line != "" {
			l.output(l.combineOutputStr(LEVEL_ERR, line))
		}
	}
}

func (l *Logger) Fatal(v ...any) {
	l.show(LEVEL_FATAL, v)

	line := lineInfo(line_info_inner_index)
	if line != "" {
		l.output(l.combineOutputStr(LEVEL_FATAL, line))
	}

	panic(v)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.showf(LEVEL_FATAL, format, v)

	line := lineInfo(line_info_inner_index)
	if line != "" {
		l.output(l.combineOutputStr(LEVEL_FATAL, line))
	}

	panic(v)
}

func (l *Logger) Fatalj(v ...any) {
	l.showj(LEVEL_FATAL, v)

	line := lineInfo(line_info_inner_index)
	if line != "" {
		l.output(l.combineOutputStr(LEVEL_FATAL, line))
	}

	panic(v)
}

func (l *Logger) Fatalji(v ...any) {
	l.showji(LEVEL_FATAL, v)

	line := lineInfo(line_info_inner_index)
	if line != "" {
		l.output(l.combineOutputStr(LEVEL_FATAL, line))
	}

	panic(v)
}

func (l *Logger) show(level int, v []any) {
	if l.level < level {
		l.tag = ""
		return
	}

	l.output(l.combineOutputStr(l.level, l.genStr(v)))
}

func (l *Logger) showf(level int, format string, v []any) {
	if l.level < level {
		l.tag = ""
		return
	}

	l.output(l.combineOutputStr(l.level, fmt.Sprintf(format, v)))
}

func (l *Logger) showj(level int, v []any) {
	if l.level < level {
		l.tag = ""
		return
	}

	l.output(l.combineOutputStr(l.level, l.genJsonStr(false, v)))
}

func (l *Logger) showji(level int, v []any) {
	if l.level < level {
		l.tag = ""
		return
	}

	l.output(l.combineOutputStr(l.level, l.genJsonStr(true, v)))
}

func (l *Logger) genStr(v []any) string {
	str := fmt.Sprint(v)
	return string([]byte(str)[1 : len(str)-1])
}

func (l *Logger) genJsonStr(indent bool, v []any) string {
	defer func() {
		if err := recover(); err != nil {
			l.Error(v...)
		}
	}()

	var slice []byte
	var err error
	newParams := make([]any, len(v))
	for index, param := range v {
		value := reflect.ValueOf(v)
		if !value.IsValid() || value.IsZero() {
			newParams[index] = " "
			continue
		}

		switch reflect.ValueOf(v).Type().Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Ptr:
			if indent {
				slice, err = json.MarshalIndent(v, "", "    ")
			} else {
				slice, err = json.Marshal(v)
			}

			if err != nil {
				newParams[index] = param
			} else {
				newParams[index] = string(slice)
			}

		default:
			newParams[index] = param
		}
	}

	return l.genStr(newParams)
}

func (l *Logger) combineOutputStr(level int, str string) string {
	tag := l.tag
	l.tag = ""

	var levelStr string
	if l.WithColor {
		levelStr = levelStrColorMap[level]
	} else {
		levelStr = levelStrMap[level]
	}

	if tag != "" {
		str = levelStr + tag + ": " + str
	} else {
		str = levelStr + str
	}

	if len(str) > 4*1024*1024 {
		str = "too long:" + strconv.Itoa(len(str)) + ",content:" + string([]byte(str)[:100]) + "..."
	}

	return str
}

func (l *Logger) output(str string) {
	switch l.kind {
	case LOG_KIND_FILE:
		if nil == l.file {
			l.createLogFile()
		}

		l.logger.Println(str)

		l.currentSize += len(log_time_fmt) + len(str) + 1
		l.checkFileRotate()

	case LOG_KIND_SIMPLE:
		l.logger.Println(str)

	case LOG_KIND_MEM:
		if nil == l.buf {
			l.buf = make([]string, l.RotateSize)
			l.bufIndex = 0
		}

		l.buf[l.bufIndex] = time.Now().Format(log_time_fmt) + " " + str
		l.bufIndex = (l.bufIndex + 1) % l.RotateSize

	case LOG_KIND_CALL:
		if nil != l.CallFunc {
			l.CallFunc(time.Now().Format(log_time_fmt), str)
		}

	default:
		log.Println(str)
	}
}

func (l *Logger) createLogFileByName(name string) {
	if nil != l.file {
		l.file.Close()
		l.file = nil
	}

	info, err := os.Stat(name)
	if err == nil {
		timeStr := info.ModTime().Format(log_file_fmt)
		newName := strings.TrimSuffix(name, log_file_ext) + "_" + timeStr + log_file_ext
		os.Rename(name, newName)
	}

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

	logFileName := fmt.Sprintf("%s/%s"+log_file_ext, l.dir, l.fileName)
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
	if l.DelHour == 0 && l.DelNum == 0 && l.DelSize == 0 {
		return
	}

	type logFileInfo struct {
		name    string
		modTime int64
		size    int64
	}

	infos := make([]*logFileInfo, 0)

	currentFileName := l.fileName + log_file_ext
	filepath.WalkDir(l.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		name := d.Name()
		if name == currentFileName {
			return nil
		}

		if !strings.HasPrefix(name, l.fileName) || !strings.HasSuffix(name, log_file_ext) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		infos = append(infos, &logFileInfo{
			name: path, modTime: info.ModTime().Unix(), size: info.Size(),
		})

		return nil
	})

	if len(infos) == 0 {
		return
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].modTime < infos[j].modTime
	})

	if l.DelHour > 0 {
		delTime := time.Now().Add(time.Duration(-l.DelHour) * time.Hour).Unix()
		for k, v := range infos {
			if v.modTime > delTime {
				infos = infos[k:]
				break
			} else {
				os.Remove(v.name)
			}
		}
	}

	if l.DelNum > 0 {
		length := len(infos)
		if length > l.DelNum {
			for _, v := range infos[:length-l.DelNum] {
				os.Remove(v.name)
			}

			infos = infos[l.DelNum:]
		}
	}

	if l.DelSize > 0 {
		var totalSize int64 = 0
		for i := len(infos) - 1; i >= 0; i-- {
			totalSize += infos[i].size
			if totalSize < l.DelSize {
				continue
			}

			os.Remove(infos[i].name)
		}
	}
}

func (l *Logger) createLoggerDir() error {
	if l.dir == "." {
		return nil
	}

	dirInfo, err := os.Stat(l.dir)
	if nil == err && dirInfo.IsDir() {
		return nil
	}

	return os.MkdirAll(l.dir, 0777)
}
