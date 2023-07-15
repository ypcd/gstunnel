package gstestpipe

import (
	"log"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

var g_logger *log.Logger

func init() {
	g_logger = gstunnellib.NewLoggerFileAndStdOut("GsTestPiPe.log")
}

func checkError(err error) {
	gstunnellib.CheckErrorEx(err, g_logger)
}

func checkError_exit(err error) {
	gstunnellib.CheckErrorEx_exit(err, g_logger)
}

func checkError_panic(err error) {
	gstunnellib.CheckErrorEx_panic(err)
}

func CheckError_test(err error, t *testing.T) {
	gstunnellib.CheckError_test(err, t)
}

func CheckError_test_noExit(err error, t *testing.T) {
	gstunnellib.CheckError_test_noExit(err, t)
}

func checkError_info(err error) {
	gstunnellib.CheckErrorEx_info(err, g_logger)
}
