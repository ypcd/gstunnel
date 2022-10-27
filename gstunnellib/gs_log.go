package gstunnellib

import (
	"io"
	"log"
	"os"
)

func init() {
	//fmt.Println("gs_log init().")
	//log.New().Output()
	//os.Stdout.Write()

}

func NewFileLogger(FileName string) *log.Logger {
	lf, err := os.Create(FileName)
	checkError(err)

	lws := io.MultiWriter(lf, os.Stdout)
	log1 := log.New(lws, "", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return log1
}

type Logger_List struct {
	GenLogger  *log.Logger
	GSIpLogger *log.Logger
	GSNetIOLen *log.Logger
}

//func NewLoggerList()
