package gstunnellib

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Benchmark_iocopy(b *testing.B) {
	testCacheSize := 16 * 1024 * 1024
	SendData := gsrand.GetRDBytes(testCacheSize)
	fd, err := os.OpenFile(os.DevNull, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	checkError_panic(err)
	defer fd.Close()

	for n := 0; n < b.N; n++ {
		_, err = io.Copy(fd, bytes.NewBuffer(SendData))
		checkError_panic(err)
	}
}

func Benchmark_NetConnWriteAll(b *testing.B) {
	testCacheSize := 16 * 1024 * 1024
	SendData := gsrand.GetRDBytes(testCacheSize)
	fd, err := os.OpenFile(os.DevNull, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	checkError_panic(err)
	defer fd.Close()

	for n := 0; n < b.N; n++ {
		_, err = netConnWriteAll_test(fd, SendData)
		checkError_panic(err)
	}
}
