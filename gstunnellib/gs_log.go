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

func CreateFileLogger(FileName string) *log.Logger {
	lf, err := os.Create(FileName)
	checkError(err)

	lws := newGs_logger_writers(lf, os.Stdout)
	logger = log.New(lws, "", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return logger
}

type gs_logger_writers struct {
	fileout io.Writer
	stdout1 io.Writer
}

func newGs_logger_writers(fileout io.Writer, stdout1 io.Writer) io.Writer {
	return io.MultiWriter(fileout, stdout1)
}

/*
func (gl *gs_logger_writers) write_old(p []byte) (int, error) {
	n, err1 := gl.fileout.Write(p)
	_, err2 := gl.stdout1.Write(p)
	var err error = err1
	if err1 != nil {
		err = err1
	} else if err2 != nil {
		err = err2
	}
	return n, err
}
*/
