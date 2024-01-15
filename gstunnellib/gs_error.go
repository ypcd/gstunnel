package gstunnellib

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

const g_gsErrorStackBufSize int = gserror.G_gsErrorStackBufSize

func checkError(err error) {
	gserror.CheckErrorEx(err, G_logger)
}

func checkError_exit(err error) {
	gserror.CheckErrorEx_exit(err, G_logger)
}

func checkError_info(err error) {
	gserror.CheckErrorEx_info(err, G_logger)
}

func checkError_panic(err error) {
	gserror.CheckErrorEx_panic(err)
}

func checkErrorEx(err error, inlogger *log.Logger) {
	gserror.CheckErrorEx(err, inlogger)
}

func checkErrorEx_panic(err error) {
	gserror.CheckErrorEx_panic(err)
}

func checkErrorEx_exit(err error, inlogger *log.Logger) {
	gserror.CheckErrorEx_exit(err, inlogger)
}

func checkErrorEx_info(err error, inlogger *log.Logger) {
	gserror.CheckErrorEx_info(err, inlogger)
}

/*
	func CheckError_test_old(inerr error, t *testing.T) {
		if inerr != nil {
			_, file, line, ok := runtime.Caller(1)
			if !ok {
				file = "???"
				line = -1
			}
			finfo, _ := os.Stat(file)
			fileName := finfo.Name()
			t.Logf("%s:%d: ", fileName, line)

			panic(inerr)
			//t.Fatal(inerr)
		}
	}

	func CheckError_test_noExit_old(inerr error, t *testing.T) {
		if inerr != nil {
			_, file, line, ok := runtime.Caller(1)
			if !ok {
				file = "???"
				line = -1
			}
			finfo, _ := os.Stat(file)
			fileName := finfo.Name()
			errstr := fmt.Sprintf("%s:%d:Test error: %s", fileName, line, inerr)
			t.Logf(errstr)
			log.Println(errstr)
		}
	}
*/
func CheckError_test(inerr error, t *testing.T)        { gserror.CheckErrorEx_panic(inerr) }
func CheckError_test_noExit(inerr error, t *testing.T) { gserror.CheckErrorEx_info(inerr, G_logger) }

// panic recover
/*
func Panic_Recover(inlog *log.Logger) {
	if msg := recover(); msg != nil {
		tmp := make([]byte, g_gsErrorStackBufSize)
		nlen := runtime.Stack(tmp, false)
		inlog.Output(2,
			fmt.Sprintf("Panic[Func exit] recover msg: %s\nPanic runtime.Stack: %s\n", msg, string(tmp[:nlen])))
	}
}
*/

func Panic_Recover_GSCtx(inlog *log.Logger, gctx IGSContext) {
	if msg := recover(); msg != nil {
		tmp := make([]byte, g_gsErrorStackBufSize)
		nlen := runtime.Stack(tmp, false)
		inlog.Output(1,
			fmt.Sprintf("Panic[Func exit] recover msg: [%d] %s\nPanic runtime.Stack: %s\n", gctx.GetGsId(), msg, string(tmp[:nlen])))
	}
}

func IsErrorNetUsually(err error) bool {
	return gserror.IsErrorNetUsually(err)
}
