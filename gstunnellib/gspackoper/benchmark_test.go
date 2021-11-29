package gspackoper

/*
goos: windows
goarch: amd64
pkg: gstunnellib
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Benchmark_bytesjoin_old-4        9936274              3631 ns/op            3960 B/op         12 allocs/op
Benchmark_bytesjoin_buf-4       11542570              3124 ns/op            1592 B/op         11 allocs/op
Benchmark_bytesjoin_pool-4      12583714              2873 ns/op             376 B/op          9 allocs/op
Benchmark_bytesjoin_gen-4       11374375              3166 ns/op            1592 B/op         11 allocs/op


goos: darwin
goarch: amd64
pkg: gstunnellib_mod
cpu: Intel(R) Core(TM) i5-8210Y CPU @ 1.60GHz
Benchmark_bytesjoin_old-2             259254          4709 ns/op        3697 B/op          6 allocs/op
Benchmark_bytesjoin_buf-2             344426          3556 ns/op        1329 B/op          5 allocs/op
Benchmark_bytesjoin_pool-2            366908          3866 ns/op         113 B/op          3 allocs/op
Benchmark_bytesjoin_gen-2             340984          6400 ns/op        1329 B/op          5 allocs/op
Benchmark_json_po1-2                   62082         17962 ns/op        1652 B/op         25 allocs/op
Benchmark_json_proto-2                110804         10527 ns/op        1016 B/op         37 allocs/op
Benchmark_json_po1_inbytes-2           24458         49005 ns/op       22242 B/op         29 allocs/op
Benchmark_json_proto_inbytes-2         65790         18357 ns/op       10425 B/op         38 allocs/op

*/

import (
	. "gstunnel/gstunnellib/gsrand"
	"testing"
)

var po1 *packOper = createPackOperChangeKey()

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

func Benchmark_json_po1(b *testing.B) {
	ap1 := NewCompresser()
	blen := 0
	for i := 0; i < b.N; i++ {
		re := jsonPacking_OperGen_po1([]byte{})
		cdata := ap1.compress(re)
		_ = cdata
		if blen == 0 {
			blen = len(cdata)
		}
	}
	//b.Log(blen)
}

func Benchmark_json_proto(b *testing.B) {
	blen := 0
	for i := 0; i < b.N; i++ {
		re := jsonPacking_OperGen([]byte{})
		if blen == 0 {
			blen = len(re)
		}
	}
	//b.Log(blen)
}

func Benchmark_json_po1_inbytes(b *testing.B) {
	ap1 := NewCompresser()
	blen := 0
	data1 := GetRDCBytes(1024 * 2)
	for i := 0; i < b.N; i++ {
		re := jsonPacking_OperGen_po1(data1)
		cdata := ap1.compress(re)
		_ = cdata
		if blen == 0 {
			blen = len(cdata)
		}
	}
	//b.Log(blen)
}

func Benchmark_json_proto_inbytes(b *testing.B) {
	blen := 0
	data1 := GetRDCBytes(1024 * 2)

	for i := 0; i < b.N; i++ {
		re := jsonPacking_OperGen(data1)
		if blen == 0 {
			blen = len(re)
		}
	}
	//b.Log(blen)
}

func Benchmark_json_un_po1_inbytes(b *testing.B) {
	ap1 := NewCompresser()
	blen := 0
	data1 := GetRDCBytes(1024 * 2)
	re := jsonPacking_OperGen_po1(data1)
	cdata := ap1.compress(re)
	for i := 0; i < b.N; i++ {
		unc := ap1.uncompress(cdata)
		_, _ = jsonUnPack_OperGen_po1(unc)
		_ = cdata
		if blen == 0 {
			blen = len(cdata)
		}
	}
	//b.Log(blen)
}

func Benchmark_json_un_proto_inbytes(b *testing.B) {
	blen := 0
	data1 := GetRDCBytes(1024 * 2)
	re := jsonPacking_OperGen(data1)

	for i := 0; i < b.N; i++ {
		unre, _ := jsonUnPack_OperGen(re)
		_ = unre
		if blen == 0 {
			blen = len(re)
		}
	}
	//b.Log(blen)
}
