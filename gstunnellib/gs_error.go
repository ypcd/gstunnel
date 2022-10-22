package gstunnellib

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"
)

// Print error info, print stack info.
func CheckErrorEx(err error, inlogger *log.Logger) {
	if err != nil {
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlogger.Output(3, fmt.Sprintf("Error: %s\nstack: %s\n", err.Error(), string(tmp[:nlen])))
	}
}

func CheckErrorEx_panic(err error, inlogger *log.Logger) {
	if err != nil {
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlogger.Output(3, fmt.Sprintf("Panic error: %s\nstack: %s\n", err.Error(), string(tmp[:nlen])))

		panic(err)
	}
}

func CheckErrorEx_exit(err error, inlogger *log.Logger) {
	if err != nil {
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlogger.Output(3, fmt.Sprintf("Fatal error: %s\nstack: %s\n", err.Error(), string(tmp[:nlen])))
		os.Exit(-1)
	}
}

func CheckErrorEx_info(err error, inlogger *log.Logger) {
	if err != nil {
		inlogger.Output(3, fmt.Sprintf("Info: %s\n", err.Error()))
	}
}

func checkError(err error) {
	CheckErrorEx(err, logger)
}

func checkError_exit(err error) {
	CheckErrorEx_exit(err, logger)
}

func checkError_info(err error) {
	CheckErrorEx_info(err, logger)
}

func checkError_panic(err error) {
	CheckErrorEx_panic(err, logger)
}

func CheckError_test(inerr error, t *testing.T) {
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

func CheckError_test_noExit(inerr error, t *testing.T) {
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

// panic recover
func Panic_Recover(inlog *log.Logger) {
	if msg := recover(); msg != nil {
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlog.Output(1,
			fmt.Sprintf("Panic[Func exit] recover msg: %s\nPanic static: %s\n", msg, string(tmp[:nlen])))
	}
}
