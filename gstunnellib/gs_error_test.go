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
	checkError(nil)
	//checkError(errors.New("123"))
}

func Test_CheckError2(t *testing.T) {
	checkError(errors.New("Error test."))
	//checkError(errors.New("123"))
}

func Test_Panic_Recover(t *testing.T) {
	defer Panic_Recover(log.Default())
	panic("Error.")
}
