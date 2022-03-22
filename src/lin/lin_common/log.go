package lin_common

import (
	"fmt"
	"os"
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
	enableLog bool
}
var globalLogMgr = LogMgr{
	chLog : make(chan *LogMsg, 1000),
	enableLog : false,
}

func LogDebug(args ... interface{}) {
	if !globalLogMgr.enableLog {
		return
	}

	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	if len(globalLogMgr.chLog) >= 999 {
		return
	}
	l := &LogMsg{}
	l.strLog = fmt.Sprintf("%s[%s:%d] route:%d %s %s\r\n",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName, fmt.Sprint(args...))
	globalLogMgr.chLog <- l
}

func LogErr(args ... interface{}) {
	if !globalLogMgr.enableLog {
		return
	}

	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	l := &LogMsg{}
	l.strLog = fmt.Sprintf("ERROR %s[%s:%d] route:%d %s %s\r\n%s\r\n",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName, fmt.Sprint(args...), fmt.Sprintf(string(debug.Stack())))
	globalLogMgr.chLog <- l
}

func InitLog(str string) {
	globalLogMgr.logFile = str
	globalLogMgr.enableLog = true
	go go_logPrint()
}

func go_logPrint() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("log err:", err)
		}
	}()

	filehandle, err := os.OpenFile(globalLogMgr.logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("log open file err:", err)
		return
	}
	count := 0
	for l := range globalLogMgr.chLog{
		if filehandle != nil {
			_, err = filehandle.WriteString(l.strLog)
		}
		if err != nil {
			fmt.Println(err)
			if filehandle != nil {
				filehandle.Close()
			}
			filehandle, err = os.OpenFile(globalLogMgr.logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println(err)
			}
		}
		count ++
		if count > 2 {
			err = filehandle.Sync()
			if err != nil {
				fmt.Println(err)
			}
			count = 0
			err = filehandle.Truncate(0)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}