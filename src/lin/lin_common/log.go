package lin_common

import (
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
	"time"
)

func LogDebug(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	fmt.Println(fmt.Sprintf("%s[%s:%d] route:%d %s",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName),
		fmt.Sprint(args...))
}

func LogErr(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	fmt.Println(fmt.Sprintf("ERROR %s[%s:%d] route:%d %s",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName),
		fmt.Sprint(args...))
	fmt.Println(string(debug.Stack()))
}

func LogDumpStack() string {
	return string(debug.Stack())
}