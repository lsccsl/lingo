package lin_common

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"time"
)

type LOG_LEVEL int
const (
	LOG_LEVEL_DEBUG LOG_LEVEL = 1
	LOG_LEVEL_INFO LOG_LEVEL = 2
	LOG_LEVEL_ERR LOG_LEVEL = 3
)

type LogMsg struct {
	strLog string
	logLevel LOG_LEVEL
}

type LogMgr struct {
	chLog chan *LogMsg
	logFile string
	logFileErr string
	enableLog bool
	enableConsolePrint bool
	enableFilePrint bool
}
var globalLogMgr = LogMgr{
	chLog : make(chan *LogMsg, 1000),
	enableLog : false,
	enableConsolePrint : false,
	enableFilePrint : true,
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
	l := &LogMsg{logLevel:LOG_LEVEL_DEBUG}
	l.strLog = fmt.Sprintf("%s[%s:%d] route:%d %s %s\r\n",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName, fmt.Sprint(args...))
	globalLogMgr.chLog <- l
}

func LogInof(args ... interface{}) {
	if !globalLogMgr.enableLog {
		return
	}

	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	if len(globalLogMgr.chLog) >= 999 {
		return
	}
	l := &LogMsg{logLevel:LOG_LEVEL_INFO}
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
	l := &LogMsg{logLevel:LOG_LEVEL_ERR}
	l.strLog = fmt.Sprintf("ERROR %s[%s:%d] route:%d %s %s\r\n%s\r\n",
		time.Now().Format(time.RFC3339Nano), path.Base(filename), line, GetGID(), funcName, fmt.Sprint(args...), fmt.Sprintf(string(debug.Stack())))
	globalLogMgr.chLog <- l
}

func InitLog(str string, strErr string, enableConsolePrint bool, enableFilePrint bool) {
	globalLogMgr.logFile = str
	globalLogMgr.logFileErr = strErr
	globalLogMgr.enableLog = true
	globalLogMgr.enableConsolePrint = enableConsolePrint
	globalLogMgr.enableFilePrint = enableFilePrint
	go go_logPrint()
}

func go_logPrint() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("log err:", err)
		}
	}()

	var filehandle *os.File = nil
	var errOpen error
	if globalLogMgr.enableFilePrint {
		filehandle, errOpen = os.OpenFile(globalLogMgr.logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if errOpen != nil {
			fmt.Println("log open file err:", errOpen)
			return
		}
	}
	var filehandleErr *os.File = nil
	if globalLogMgr.enableFilePrint {
		filehandleErr, errOpen = os.OpenFile(globalLogMgr.logFileErr, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if errOpen != nil {
			fmt.Println("log open file err:", errOpen)
			return
		}
	}

	count := 0
	for l := range globalLogMgr.chLog{
		if globalLogMgr.enableConsolePrint {
			fmt.Print(l.strLog)
		}
		if globalLogMgr.enableFilePrint{
			if filehandle != nil {
				_, err := filehandle.WriteString(l.strLog)
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
			}

			if l.logLevel == LOG_LEVEL_ERR {
				if filehandleErr != nil {
					filehandleErr.WriteString(l.strLog)
				}
			}
		}
		count ++
		if count > 20000 {
			err := filehandle.Sync()
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