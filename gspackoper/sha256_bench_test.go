package gspackoper

import (
	"crypto/sha256"
	"log"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_sha256_bench(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 100; i++ {
		sha256.Sum256(gsrand.GetRDBytes(4 * 1024))
	}
	t2 := time.Since(t1)
	log.Println("sha256[4*1024] run time:", t2.Seconds())

	t1 = time.Now()
	sha256.Sum256(gsrand.GetRDBytes(100 * 4 * 1024))
	t2 = time.Since(t1)
	log.Println("sha256[10*4*1024] run time:", t2.Seconds())

}
