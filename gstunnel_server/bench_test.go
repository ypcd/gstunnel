package main

/*
goos: windows
goarch: amd64
pkg: github.com/ypcd/gstunnel/v6/gstunnel_server
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Benchmark_ioCopy-4         	122449850	        10.09 ns/op	       1 B/op	       0 allocs/op
Benchmark_ioCopyBuffer-4   	99957681	        11.90 ns/op	       2 B/op	       0 allocs/op
Benchmark_NoIoCopy-4       	708474978	         1.760 ns/op	       0 B/op	       0 allocs/op

goos: windows
goarch: amd64
pkg: github.com/ypcd/gstunnel/v6/gstunnel_server
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Benchmark_ioCopy-4         	128825271	         9.025 ns/op	       1 B/op	       0 allocs/op
Benchmark_ioCopyBuffer-4   	116722004	        10.19 ns/op	       	   1 B/op	       0 allocs/op
Benchmark_NoIoCopy-4       	686777587	         1.652 ns/op	       0 B/op	       0 allocs/op
*/

import (
	"bytes"
	"io"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Benchmark_ioCopy(b *testing.B) {
	src := bytes.NewBuffer(gsrand.GetRDBytes(100 * 1024 * 1024))
	dst := bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		_, err := io.Copy(&dst, src)
		checkError(err)
	}

}

func Benchmark_ioCopyBuffer(b *testing.B) {
	src := bytes.NewBuffer(gsrand.GetRDBytes(100 * 1024 * 1024))
	dst := bytes.Buffer{}
	buf := make([]byte, 4*1024)

	for i := 0; i < b.N; i++ {
		_, err := io.CopyBuffer(&dst, src, buf)
		checkError(err)
	}
}

func Benchmark_NoIoCopy(b *testing.B) {
	src := bytes.NewBuffer(gsrand.GetRDBytes(100 * 1024 * 1024))
	dst := bytes.Buffer{}
	buf := make([]byte, 4*1024)

	for i := 0; i < b.N; i++ {
		for {
			n, err := src.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			checkError(err)
			dst.Write(buf[:n])
		}
	}

}
