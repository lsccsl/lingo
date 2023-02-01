package common

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
)

type ERR_NUM int

const (
	ERR_NONE               ERR_NUM = 0
	ERR_sys                ERR_NUM = 1
	ERR_not_tcp_connection ERR_NUM = 2
	ERR_rpc_timeout        ERR_NUM = 3
	ERR_no_srv             ERR_NUM = 4
	ERR_rpc_response_err   ERR_NUM = 5
)

type MyError struct {
	Errfile string
	Errline int
	ErrFunc   string
	ErrNo     ERR_NUM
	ErrString string
}

func (errNo ERR_NUM)String() string {
	errName := "SRV_TYPE_none"
	switch errNo {
	case ERR_sys:
		errName = "ERR_sys"
	case ERR_not_tcp_connection:
		errName = "ERR_not_tcp_connection"
	case ERR_rpc_timeout:
		errName = "ERR_rpc_timeout"
	case ERR_no_srv:
		errName = "ERR_no_srv"
	case ERR_rpc_response_err:
		errName = "ERR_rpc_response_err"
	}
	return "[err:" + errName + "(" + strconv.FormatInt(int64(errNo), 10) + ")]"
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

	err.ErrString = fmt.Sprintf("[%s:%d]%s", path.Base(filename), line, funcName) +
		errNo.String() +
		fmt.Sprint(param...)

	return err
}

func GenErrNoERR_NUM(param... interface{})*MyError {
	pc,filename, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	err := &MyError{
		Errfile: path.Base(filename),
		Errline: line,
		ErrFunc: funcName,
		ErrNo:   ERR_NONE,
	}

	err.ErrString = fmt.Sprintf("[%s:%d]%s", path.Base(filename), line, funcName) + fmt.Sprint(param...)

	return err
}