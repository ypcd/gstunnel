package gstunnellib

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func init() {
	//fmt.Println("gs_log init().")
	//log.New().Output()
	//os.Stdout.Write()

}

func GetNowTimeString() string {
	t1 := time.Now()
	return fmt.Sprintf("%d-%d-%d_%d-%d-%d", t1.Year(), t1.Month(), t1.Day(),
		t1.Hour(), t1.Minute(), t1.Second())
}

func GetFileNameAddTime(FileName string) string {

	err := os.Mkdir("logs", 0774)
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			panic(err)
		}
	}
	FileName = "logs/" + FileName

	rn := strings.IndexAny(FileName, ".")
	if rn == -1 {
		return FileName + "_" + GetNowTimeString() + "_"
	} else {
		list1 := strings.Split(FileName, ".")
		return list1[0] + "_" + GetNowTimeString() + "_." + strings.Join(list1[1:], ".")
	}
	//return FileName +"_" GetNowTimeString()
}

func GetExeDir() string {
	s1 := os.Args[0]
	s2 := strings.ReplaceAll(s1, "\\", "/")
	d1 := path.Dir(s2)
	return d1
}

func newLoggerFileAndStdOutEx(FileName string, outFileName *string) *log.Logger {
	var filePath string = GetFileNameAddTime(FileName)
	if outFileName != nil {
		*outFileName = filePath
	}
	/*
		if !strings.Contains(FileName, "\\") && !strings.Contains(FileName, "/") {
			filePath = path.Join(GetExeDir(), FileName)
		}

			fmt.Println("args:", os.Args)
			fmt.Println("filePath:", filePath)
			dir, _ := os.Getwd()
			fmt.Println("word dir:", dir)
	*/
	lf, err := os.Create(filePath)
	CheckError_exit(err)

	lws := io.MultiWriter(lf, os.Stdout)
	log1 := log.New(lws, "", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return log1
}

func NewLoggerFileAndStdOut(FileName string) *log.Logger {
	return newLoggerFileAndStdOutEx(FileName, nil)
}

func NewLoggerFile(FileName string) *log.Logger {
	var filePath string = GetFileNameAddTime(FileName)
	/*
		if !strings.Contains(FileName, "\\") && !strings.Contains(FileName, "/") {
			filePath = path.Join(GetExeDir(), FileName)
		}

			fmt.Println("args:", os.Args)
			fmt.Println("filePath:", filePath)
			dir, _ := os.Getwd()
			fmt.Println("word dir:", dir)
	*/
	lf, err := os.Create(filePath)
	CheckError_exit(err)

	log1 := log.New(lf, "", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return log1
}

func NewLoggerFileAndLog(FileName string, inlog io.Writer) *log.Logger {
	var filePath string = GetFileNameAddTime(FileName)
	/*
		if !strings.Contains(FileName, "\\") && !strings.Contains(FileName, "/") {
			filePath = path.Join(GetExeDir(), FileName)
		}

			fmt.Println("args:", os.Args)
			fmt.Println("filePath:", filePath)
			dir, _ := os.Getwd()
			fmt.Println("word dir:", dir)
	*/
	lf, err := os.Create(filePath)
	CheckError_exit(err)

	mw := io.MultiWriter(lf, inlog)
	log1 := log.New(mw, "", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return log1
}

type Logger_List struct {
	GenLogger  *log.Logger
	GSIpLogger *log.Logger
	GSNetIOLen *log.Logger
}
