package gstunnellib

import "testing"

func Test_checkError_test(t *testing.T) {
	checkError_test(nil, t)
	//checkError_test(errors.New("123"), t)
}

func Test_CheckError(t *testing.T) {
	checkError(nil)
	//checkError(errors.New("123"))
}
