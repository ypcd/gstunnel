package gsmap

import (
	"strconv"
	"sync"
	"testing"
)

func Test_gsmap_add(t *testing.T) {
	map1 := NewGSMap()

	for i := 0; i < 100; i++ {
		map1.Add(strconv.Itoa(i), strconv.Itoa(i))
	}

	var ks, vs int
	for k, v := range map1.data {
		k1, _ := strconv.Atoi(k)
		v1, _ := strconv.Atoi(v)
		ks += k1
		vs += v1
	}
	if ks != vs || ks != 99*50 {
		panic("ks!=vs")
	}
}

func Test_gsmap_add_mt(t *testing.T) {
	map1 := NewGSMap()

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			map1.Add(strconv.Itoa(v), strconv.Itoa(v))
		}(i)
	}

	wg.Wait()

	var ks, vs int
	for k, v := range map1.data {
		k1, _ := strconv.Atoi(k)
		v1, _ := strconv.Atoi(v)
		ks += k1
		vs += v1
	}
	if ks != vs || ks != 99*50 {
		panic("ks!=vs")
	}
}
