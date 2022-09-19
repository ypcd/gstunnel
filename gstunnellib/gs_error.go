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
		//fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		inlogger.Output(3, fmt.Sprintln("Error:", err.Error()))
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlogger.Println("Error stack:", string(tmp[:nlen]))
	}
}

func CheckErrorEx_panic(err error, inlogger *log.Logger) {
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		inlogger.Output(3, fmt.Sprintln("Fatal error:", err.Error()))
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlogger.Println("Fatal error stack:", string(tmp[:nlen]))
		panic(err)
	}
}

func CheckErrorEx_exit(err error, inlogger *log.Logger) {
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		inlogger.Output(3, fmt.Sprintln("Fatal error:", err.Error()))
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlogger.Println("Fatal error stack:", string(tmp[:nlen]))
		os.Exit(-1)
	}
}

func checkError(err error) {
	CheckErrorEx(err, logger)
}

func checkError_exit(err error) {
	CheckErrorEx_exit(err, logger)
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
		t.Logf("%s:%d: %s", fileName, line, inerr)
		log.Println("Error:", inerr)
	}
}

// panic recover
func Panic_Recover(inlog *log.Logger) {
	if x := recover(); x != nil {
		inlog.Println("Panic[Func exit] recover msg:", x)
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, false)
		inlog.Println("Panic stack:")
		inlog.Println(string(tmp[:nlen]))
	}
}
