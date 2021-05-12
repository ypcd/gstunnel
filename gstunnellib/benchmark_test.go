package gstunnellib

/*
goos: windows
goarch: amd64
pkg: gstunnellib
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Benchmark_bytesjoin_old-4        9936274              3631 ns/op            3960 B/op         12 allocs/op
Benchmark_bytesjoin_buf-4       11542570              3124 ns/op            1592 B/op         11 allocs/op
Benchmark_bytesjoin_pool-4      12583714              2873 ns/op             376 B/op          9 allocs/op
Benchmark_bytesjoin_gen-4       11374375              3166 ns/op            1592 B/op         11 allocs/op

*/

import (
	"testing"
)

var po1 *packOper = CreatePackOperChangeKey()

func init() {
	po1.Data = GetRDBytes(1024)
}

func Benchmark_bytesjoin_old(b *testing.B) {
	for i := 0; i < b.N; i++ {
		po1.GetSha256_old()
	}
}

func Benchmark_bytesjoin_buf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		po1.GetSha256_buf()
	}
}

func Benchmark_bytesjoin_pool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		po1.GetSha256_pool()
	}
	/*
		for i := 0; i < 100; i++ {
			po1.GetSha256_pool()
		}
	*/
}

func Benchmark_bytesjoin_gen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		po1.GetSha256()
	}
}
