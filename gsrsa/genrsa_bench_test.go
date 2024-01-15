package gsrsa

import (
	"testing"
)

func Benchmark_genrsa(b *testing.B) {

	for n := 0; n < b.N; n++ {
		key := NewGenRSAObj(g_rsaKeyLenbits)
		_ = key
	}
}

func Benchmark_genrsa2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		key, key2 := GenerateKeyPair(g_rsaKeyLenbits)
		_, _ = key, key2
	}
}
