package gstunnellib

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var gslogtest_outFileName string

func Test_gs_log(t *testing.T) {
	outf, err := os.Create("tmp.log")
	CheckError_test(err, t)

	lws := io.MultiWriter(outf, os.Stdout)
	g_logger := log.New(lws, "", log.LstdFlags|log.Lshortfile)
	g_logger.Println("test.")
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
	var outFileName *string = &gslogtest_outFileName
	lg := newLoggerFileAndStdOutEx("tmp.log", outFileName)
	lg.Println("Log test.")
	time.Sleep(time.Second)

	inf, err := os.Open(*outFileName)
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

/*
	func Test_NewFileLogger_remove(t *testing.T) {
		err := os.Remove(gslogtest_outFileName)
		CheckError_test(err, t)
	}
*/
func Test_FileNameAddTime1(t *testing.T) {
	str1 := GetFileNameAddTime("123")
	_ = str1
	str2 := GetFileNameAddTime("123.log")
	str3 := GetFileNameAddTime("123.tmp.log")
	_, _ = str2, str3
}
