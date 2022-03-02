package lin_common

import (
	"fmt"
	"path"
	"runtime"
)

const (
	ERR_NONE = 0;
	ERR_no_dialData = 1;
)

type MyError struct {
	Errfile string
	Errline int
	ErrFunc string
	ErrNo int
	ErrString string
}

func (pthis*MyError)Error() string {
	return pthis.ErrString
}
func GenErr(errNo int, param... interface{})*MyError {
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
