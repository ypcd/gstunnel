package gslog

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
	checkError_panic(err)

	lws := io.MultiWriter(outf, os.Stdout)
	G_logger := log.New(lws, "", log.LstdFlags|log.Lshortfile)
	G_logger.Println("test.")
	inf, err := os.Open("tmp.log")
	checkError_panic(err)
	defer func() {
		outf.Close()
		inf.Close()
		err := os.Remove("tmp.log")
		checkError_panic(err)
	}()

	re, err := io.ReadAll(inf)
	checkError_panic(err)

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
	checkError_panic(err)
	defer func() {
		inf.Close()
	}()

	re, err := io.ReadAll(inf)
	checkError_panic(err)

	if !strings.Contains(string(re), "Log test.") {
		t.Fatal(err)
	}
}

func _Test_NewFileLogger_nil_g_logger(t *testing.T) {
	logger2 := G_logger
	G_logger = nil
	defer func() {
		G_logger = logger2
	}()

	var outFileName *string = &gslogtest_outFileName

	lg := NewLoggerFileAndStdOut(":")
	lg.Println("Log test.")
	time.Sleep(time.Second)

	inf, err := os.Open(*outFileName)
	checkError_panic(err)
	defer func() {
		inf.Close()
	}()

	re, err := io.ReadAll(inf)
	checkError_panic(err)

	if !strings.Contains(string(re), "Log test.") {
		t.Fatal(err)
	}
}

/*
	func Test_NewFileLogger_remove(t *testing.T) {
		err := os.Remove(gslogtest_outFileName)
		checkError_panic(err)
	}
*/
func Test_FileNameAddTime1(t *testing.T) {
	str1 := GetFileNameAddTime("123")
	_ = str1
	str2 := GetFileNameAddTime("123.log")
	str3 := GetFileNameAddTime("123.tmp.log")
	_, _ = str2, str3
}
