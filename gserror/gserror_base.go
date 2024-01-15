package gserror

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gslog"
)

const G_gsErrorStackBufSize int = 2 * 1024 * 1024

var g_logger *log.Logger

func init() {
	g_logger = gslog.NewLoggerFileAndStdOut("gs_error_base.log")
}

// Print error info, print stack info.
func CheckErrorEx(err error, inlogger *log.Logger) {
	if err != nil {
		if inlogger != nil {
			tmp := make([]byte, G_gsErrorStackBufSize)
			nlen := runtime.Stack(tmp, false)
			inlogger.Output(3, fmt.Sprintf("Error: %s\nstack: %s\n", err.Error(), string(tmp[:nlen])))
		} else {
			tmp := make([]byte, G_gsErrorStackBufSize)
			nlen := runtime.Stack(tmp, false)
			fmt.Printf("Error: %s\nstack: %s\n", err.Error(), string(tmp[:nlen]))
			//panic(err)
		}
	}
}

func CheckErrorEx_panic(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckError_panic(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckErrorEx_exit(err error, inlogger *log.Logger) {
	if err != nil {
		CheckErrorEx(err, inlogger)
		os.Exit(-1)
	}
}

func CheckErrorEx_info(err error, inlogger *log.Logger) {
	if err != nil {
		inlogger.Output(3, fmt.Sprintf("Info: %s\n", err.Error()))
	}
}

func CheckError_info(err error) {
	CheckErrorEx_info(err, g_logger)
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
func CheckError_test(inerr error, t *testing.T) { CheckErrorEx_panic(inerr) }

// panic recover
func Panic_Recover(inlog *log.Logger) {
	if msg := recover(); msg != nil {
		tmp := make([]byte, G_gsErrorStackBufSize)
		nlen := runtime.Stack(tmp, false)
		inlog.Output(2,
			fmt.Sprintf("Panic[Func exit] recover msg: %s\nPanic runtime.Stack: %s\n", msg, string(tmp[:nlen])))
	}
}

/*
func IsErrorNetUsually(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, net.ErrClosed) || errors.Is(err, os.ErrDeadlineExceeded) ||
		errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE) || errors.Is(err, io.ErrClosedPipe) {
		return true
	}
	return false
}
*/
