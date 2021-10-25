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
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		inlogger.Output(2, fmt.Sprintln("Fatal error:", err.Error()))
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
