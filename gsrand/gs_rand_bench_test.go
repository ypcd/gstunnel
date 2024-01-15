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


goos: darwin
goarch: arm64
pkg: github.com/ypcd/gstunnel/v6/gstunnellib/gsrand
Benchmark_randc_Reader_Read-8   	 2450156	       464.7 ns/op	      16 B/op	       1 allocs/op
Benchmark_randc_int-8           	  198339	      5957 ns/op	     784 B/op	      49 allocs/op
Benchmark_randc_f64-8           	 3172924	       379.0 ns/op	      48 B/op	       3 allocs/op
Benchmark_randc_f64_2-8         	 2851318	       420.7 ns/op	      63 B/op	       4 allocs/op
*/

// 新的GetRDBytes函数的性能是getRDBytes_old函数的22.4倍
func Benchmark_randc_Reader_Read(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rd := GetRDBytes(16)
		_ = rd
	}
}

func Benchmark_randc_int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rd := getRDBytes_old(16)
		_ = rd
	}
}

func Benchmark_randc_f64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rd := GetRDF64()
		_ = rd
	}
}

/*
func Benchmark_randc_f64_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rd := getRDF64_2()
		_ = rd
	}
}
*/
