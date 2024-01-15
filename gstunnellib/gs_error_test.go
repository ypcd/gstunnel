package gstunnellib

import (
	"errors"
	"testing"
	"time"
)

/*
	func Test_checkError_test(t *testing.T) {
		checkError_(nil)
		//checkError_test(errors.New("123"), t)
	}

	func Test_checkError(t *testing.T) {
		checkError(nil)
		//checkError(errors.New("123"))
	}
*/
func Test_checkError2(t *testing.T) {
	checkError(errors.New("---error test, not error---Error test."))
	//checkError(errors.New("123"))
}

func Test_checkError_exit(t *testing.T) {
	//checkError_exit(errors.New("Error test."))
	//checkError(errors.New("123"))
}

func Test_checkError_info(t *testing.T) {
	checkError_info(errors.New("Hello."))
	//checkError(errors.New("123"))
}

func noTest_panic(t *testing.T) {
	checkError_panic(errors.New("error."))
}

func Test_chan_bad(t *testing.T) {
	c1, c2 := GetRDNetConn()
	c1.Close()
	c2.Close()
	time.Sleep(time.Millisecond * 10)
	buf := make([]byte, 100)
	n1, err1 := c2.Read(buf)
	checkErrorEx(err1, G_logger)
	n2, err2 := c2.Write(make([]byte, 100))
	_, _ = n1, n2
	checkError(err2)
}
