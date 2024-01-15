package gserror

import (
	"errors"
	"log"
	"testing"
)

func Test_CheckError_test(t *testing.T) {
	CheckError_test(nil, t)
	//CheckError_test(errors.New("123"), t)
}

func Test_CheckError(t *testing.T) {
	CheckErrorEx(nil, log.Default())
	//CheckError(errors.New("123"))
}

func Test_CheckError2(t *testing.T) {
	CheckErrorEx(errors.New("---error test, not error---Error test."), log.Default())
	//CheckError(errors.New("123"))
}

func Test_CheckError3(t *testing.T) {
	CheckErrorEx(errors.New("---error test, not error---Error test."), nil)
	//CheckError(errors.New("123"))
}

func Test_CheckError_exit(t *testing.T) {
	//CheckError_exit(errors.New("Error test."))
	//CheckError(errors.New("123"))
}

func Test_Panic_Recover(t *testing.T) {
	defer Panic_Recover(log.Default())
	panic("---error test, not error---An exception occurred.")
}

func Test_CheckError_info(t *testing.T) {
	CheckErrorEx_info(errors.New("Hello."), log.Default())
	//CheckError(errors.New("123"))
}

func noTest_panic(t *testing.T) {
	CheckErrorEx_panic(errors.New("error."))
}

func Test_panic2(t *testing.T) {
	defer Panic_Recover(log.Default())
	CheckErrorEx_panic(errors.New("error."))
}
