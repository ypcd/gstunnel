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
	CheckError_test(err, t)

	lws := io.MultiWriter(outf, os.Stdout)
	logger := log.New(lws, "", log.LstdFlags|log.Lshortfile)
	logger.Println("test.")
	inf, err := os.Open("tmp.log")
	CheckError_test(err, t)
	defer func() {
		outf.Close()
		inf.Close()
		err := os.Remove("tmp.log")
		CheckError_test(err, t)
	}()

	re, err := io.ReadAll(inf)
	CheckError_test(err, t)

	if !strings.Contains(string(re), "test.") {
		t.Fatal(err)
	}
}

func Test_NewFileLogger(t *testing.T) {
	lg := NewFileLogger("tmp.log")
	lg.Println("Log test.")

	inf, err := os.Open("tmp.log")
	CheckError_test(err, t)
	defer func() {
		inf.Close()
	}()

	re, err := io.ReadAll(inf)
	CheckError_test(err, t)

	if !strings.Contains(string(re), "Log test.") {
		t.Fatal(err)
	}
}
