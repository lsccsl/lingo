package log

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

func LogDebug(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	fmt.Println(fmt.Sprintf("%s[%s:%d]%s",
		time.Now().Format(time.StampNano),
		path.Base(filename), line, funcName),
		fmt.Sprint(args...))
}

func LogErr(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	fmt.Println(fmt.Sprintf("%s[%s:%d]%s",
		time.Now().Format(time.StampNano),
		path.Base(filename), line, funcName),
		fmt.Sprint(args...))
}