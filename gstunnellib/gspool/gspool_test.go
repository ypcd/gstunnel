package gspool

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_pool_map(t *testing.T) {
	pl := NewPoolMap(8192, 1000)
	wg := sync.WaitGroup{}
	loop_total := 8000
	for i := 0; i < loop_total; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(400)))
			re := pl.Get()
			re[count%(loop_total/256+1)]++
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(400)))
			pl.Put(re)
		}(i)
	}
	wg.Wait()
	fmt.Println("size:", pl.Size())

	pl.ClearAll()
	pl.print()

	/*
		tt := 0
		sz := pl.Size()

			for i := 0; i < sz; i++ {
				for _, v := range pl.Get() {
					tt += int(v)
				}
				//fmt.Println("for tt:", tt)
			}
			fmt.Println("size:", pl.Size())
			fmt.Println("tt:", tt)
			pl.print()
	*/
}

func Test_no_pool_map(t *testing.T) {
	wg := sync.WaitGroup{}
	m1 := make([](*[]byte), 8000)
	for i := 0; i < 8000; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			re := make([]byte, 8192)
			for i := range re {
				re[i]++
			}
			m1[count] = &re
			//re[count%5000]++
			//time.Sleep(time.Millisecond * time.Duration(rand.Int31n(400)))
		}(i)
	}
	wg.Wait()
	vx := 123
	_ = vx
}
