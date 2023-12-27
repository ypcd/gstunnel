package gsrand

import (
	"testing"
)

/*
goos: windows
goarch: amd64
pkg: github.com/ypcd/gstunnel/v6/gstunnellib/gsrand
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Benchmark_randc_Reader_Read-4   	 7525699	       164.0 ns/op	      40 B/op	       2 allocs/op
Benchmark_randc_int-4           	  303093	      3679 ns/op	    1169 B/op	      65 allocs/op
*/

// 新的GetRDCBytes函数的性能是getRDCBytes_old函数的22.4倍
func Benchmark_randc_Reader_Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rd := GetRDCBytes(16)
		_ = rd
	}
}

func Benchmark_randc_int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rd := getRDCBytes_old(16)
		_ = rd
	}
}
