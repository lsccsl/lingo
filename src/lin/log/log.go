package log

import (
	"fmt"
	"lin/lin_common"
	"path"
	"runtime"
	"time"
)

func LogDebug(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	fmt.Println(fmt.Sprintf("%s[%s:%d] route:%d %s",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, lin_common.GetGID(), funcName),
		fmt.Sprint(args...))
}

func LogErr(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	fmt.Println(fmt.Sprintf("%s[%s:%d] route:%d %s",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, lin_common.GetGID(), funcName),
		fmt.Sprint(args...))
}