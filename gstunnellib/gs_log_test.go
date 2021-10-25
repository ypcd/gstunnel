package gstunnellib

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_gs_log(t *testing.T) {
	outf, err := os.Create("tmp.log")
	checkError_test(err, t)

	lws := newGs_logger_writers(outf, os.Stdout)
	logger := log.New(lws, "", log.LstdFlags|log.Lshortfile)
	logger.Println("test.")
	inf, err := os.Open("tmp.log")
	checkError_test(err, t)

	re, err := io.ReadAll(inf)
	checkError_test(err, t)

	if !strings.Contains(string(re), "test.") {
		t.Fatal(err)
	}
}

func Test_CreateFileLogger(t *testing.T) {
	lg := CreateFileLogger("tmp.log")
	lg.Println("Log test.")

	inf, err := os.Open("tmp.log")
	checkError_test(err, t)

	re, err := io.ReadAll(inf)
	checkError_test(err, t)

	if !strings.Contains(string(re), "Log test.") {
		t.Fatal(err)
	}
}
