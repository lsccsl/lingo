package lin_common

import (
	"fmt"
	"path"
	"runtime"
)

type ERR_NUM int

const (
	ERR_NONE ERR_NUM = 0
	ERR_no_dialData ERR_NUM = 1
	ERR_not_tcp_connection ERR_NUM = 2
	ERR_need_not_dial ERR_NUM = 3
)

type MyError struct {
	Errfile string
	Errline int
	ErrFunc string
	ErrNo ERR_NUM
	ErrString string
}

func (pthis*MyError)Error() string {
	return pthis.ErrString
}
func GenErr(errNo ERR_NUM, param... interface{})*MyError {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	err := &MyError{
		Errfile:path.Base(filename),
		Errline:line,
		ErrFunc:funcName,
		ErrNo:errNo,
	}

	err.ErrString = fmt.Sprintf("[%s:%d]%s", path.Base(filename), line, funcName) + fmt.Sprint(param...)

	return err
}
