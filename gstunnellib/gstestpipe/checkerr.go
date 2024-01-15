package gstestpipe

import (
	"log"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gslog"
)

var g_logger *log.Logger

func init() {
	g_logger = gslog.NewLoggerFileAndStdOut("GsTestPiPe.log")
}

func checkError(err error) {
	gserror.CheckErrorEx(err, g_logger)
}

func checkError_exit(err error) {
	gserror.CheckErrorEx_exit(err, g_logger)
}

func checkError_panic(err error) {
	gserror.CheckErrorEx_panic(err)
}

func CheckError_test(err error, t *testing.T) {
	gserror.CheckError_test(err, t)
}

/*
	func CheckError_test_noExit(err error, t *testing.T) {
		gserror.checkError_info(err)
	}
*/
func checkError_info(err error) {
	gserror.CheckErrorEx_info(err, g_logger)
}
