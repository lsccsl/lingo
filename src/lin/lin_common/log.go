package lin_common

import (
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
	"time"
)

type LogMsg struct {
	strLog string
}

type LogMgr struct {
	chLog chan *LogMsg
	logFile string
}
var globalLogMgr = LogMgr{
	chLog : make(chan *LogMsg, 1000),
}

func LogDebug(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	if len(globalLogMgr.chLog) >= 1000 {
		return
	}

	l := &LogMsg{}
	l.strLog = fmt.Sprintf(fmt.Sprintf("%s[%s:%d] route:%d %s\r\n",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName),
		fmt.Sprint(args...))

	globalLogMgr.chLog <- l
}

func LogErr(args ... interface{}) {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	l := &LogMsg{}
	l.strLog = fmt.Sprintf(fmt.Sprintf("ERROR %s[%s:%d] route:%d %s\r\n%s\r\n",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName),
		fmt.Sprint(args...),
		fmt.Sprintf(string(debug.Stack())))

	globalLogMgr.chLog <- l
}

func InitLog(str string) {
	globalLogMgr.logFile = str
	go go_logPrint()
}

func go_logPrint() {
	for l := range globalLogMgr.chLog{
		fmt.Println(l.strLog)
	}
}