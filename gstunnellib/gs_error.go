package gstunnellib

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"
)

func CheckErrorEx(err error, inlogger *log.Logger) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		inlogger.Output(3, fmt.Sprintln("Fatal error:", err.Error()))
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, true)
		inlogger.Println("Fatal error stack:", string(tmp[:nlen]))
		os.Exit(-1)
	}
}

func checkError(err error) {
	CheckErrorEx(err, logger)
}

func checkError_test(inerr error, t *testing.T) {
	if inerr != nil {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = -1
		}
		finfo, _ := os.Stat(file)
		fileName := finfo.Name()
		t.Logf("%s:%d: ", fileName, line)

		t.Fatal(inerr)
	}
}

func Panic_exit(inlog *log.Logger) {
	if x := recover(); x != nil {
		inlog.Println("Panic:", x, "  Go exit.")
		tmp := make([]byte, 1024*1024)
		nlen := runtime.Stack(tmp, true)
		inlog.Println("Panic stack:", string(tmp[:nlen]))

	}
}
