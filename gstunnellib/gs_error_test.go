package gstunnellib

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
	CheckError(nil)
	//CheckError(errors.New("123"))
}

func Test_CheckError2(t *testing.T) {
	CheckError(errors.New("---error test, not error---Error test."))
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
	CheckError_info(errors.New("Hello."))
	//CheckError(errors.New("123"))
}

func noTest_panic(t *testing.T) {
	CheckError_panic(errors.New("error."))
}
