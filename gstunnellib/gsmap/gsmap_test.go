package gsmap

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func inTest_gsmap_str_str_add(map1 *GsMapStr_Str, begin, end int) {
	//map1 := NewGSMapStr_Str()

	keyls := make([]string, 0, 100)
	for i := begin; i < end; i++ {
		map1.Add(strconv.Itoa(i), strconv.Itoa(i))
		keyls = append(keyls, strconv.Itoa(i))
	}

	//var ks, vs int
	for _, k := range keyls {
		if k != map1.Get(k) {
			panic("k!=v")
		}
	}

}

func Test_gsmap_add(t *testing.T) {
	map1 := NewGSMapStr_Str()

	inTest_gsmap_str_str_add(map1, 0, 10000)

}
func Test_gsmap_add_mt(t *testing.T) {
	map1 := NewGSMapStr_Str()

	wg := sync.WaitGroup{}
	for i := 0; i < 16; i++ {
		count := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			inTest_gsmap_str_str_add(map1, count*10000+0, count*10000+10000)
		}()
	}
	time.Sleep(time.Millisecond * 10)
	wg.Wait()
}

func Test_gsmap_str_pubkey_get_errorkey(t *testing.T) {
	map1 := NewGSMapStr_RSA()

	_, ok := map1.Get("123")
	if ok {
		panic("error")
	}
}

func inTest_gsmap_str_pubkey_add(map1 *GsMapStr_RSA, begin, end int) {
	//map1 := NewGSMapStr_Str()
	//_, pub := gsrsa.GenerateKeyPair(256)

	//keyls := make([]string, 0, 100)
	//vls := make([]*gsrsa.RSA, 0, 100)

	mapls := make(map[string]*gsrsa.RSA)

	for i := begin; i < end; i++ {
		_, pub := gsrsa.GenerateKeyPair(256)
		r1 := gsrsa.NewRSAObjPub(pub)
		k := strconv.Itoa(i)
		map1.Add(k, r1)
		mapls[k] = r1
	}

	//var ks, vs int
	for k, v := range mapls {
		rv, _ := map1.Get(k)
		if v != rv {
			panic("k!=v")
		}
	}

}

func Test_gsmap_str_pubkey_add(t *testing.T) {
	map1 := NewGSMapStr_RSA()

	inTest_gsmap_str_pubkey_add(map1, 0, 1000)

}

func Test_gsmap_str_pubkey_add_mt(t *testing.T) {
	map1 := NewGSMapStr_RSA()

	num := 100

	wg := sync.WaitGroup{}
	for i := 0; i < 16; i++ {
		count := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			inTest_gsmap_str_pubkey_add(map1, count*num+0, count*num+num)
		}()
	}
	time.Sleep(time.Millisecond * 10)
	wg.Wait()
}
