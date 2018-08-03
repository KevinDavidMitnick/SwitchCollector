package g

import (
	"encoding/json"
	"runtime"

	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type MyJSONFormatter struct {
	Time  string `json:"time"`
	File  string `json:"file"`
	Line  int    `json:"line"`
	Level string `json:"level"`
	Info  string `json:"info"`
	Msg   string `json:"msg"`
}

func (f *MyJSONFormatter) Format(entry *log.Entry) ([]byte, error) {

	logrusJF := &(log.JSONFormatter{})
	bytes, _ := logrusJF.Format(entry)

	json.Unmarshal(bytes, &f)
	if _, file, no, ok := runtime.Caller(8); ok {
		f.File = file
		f.Line = no
	}

	index := strings.Index(f.Time, "+")
	times := strings.Replace(f.Time[0:index], "T", " ", 1)
	str := fmt.Sprintf("[%s] %s %s:%d %s\n", f.Level, times, f.File, f.Line, f.Msg)
	return []byte(str), nil
}

func init() {
	log.SetFormatter(&MyJSONFormatter{})
}
