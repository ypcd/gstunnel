package gstunnellib

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

func Test_gid1(t *testing.T) {
	gid := NewGIdImp()
	wg := sync.WaitGroup{}
	chan1 := make(chan []uint64, 100)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			list1 := []uint64{}
			for i := 0; i < 1000; i++ {
				list1 = append(list1, gid.GenerateId())
			}
			chan1 <- list1
		}()
	}
	wg.Wait()
	close(chan1)

	for list1 := range chan1 {
		fmt.Println(len(list1), list1)
	}

}

func Test_gid2(t *testing.T) {
	gid := NewGIdImp()
	wg := sync.WaitGroup{}

	gon := uint64(1000)
	idn := uint64(10000)
	sumid := (gon*idn + 1) * (gon * idn / 2)

	chan1 := make(chan []uint64, gon)

	for i := uint64(0); i < gon; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			list1 := []uint64{}
			for i := uint64(0); i < idn; i++ {
				list1 = append(list1, gid.GenerateId())
			}
			chan1 <- list1
		}()
	}
	wg.Wait()
	close(chan1)

	lists := []uint64{}
	for list1 := range chan1 {
		//fmt.Println(len(list1), list1)
		lists = append(lists, list1...)
	}

	sort.Slice(lists, func(i, j int) bool { return lists[i] < lists[j] })

	sum := uint64(0)
	for _, v := range lists {
		sum += v
	}
	if sum != sumid {
		panic("Error.")
	}

	for i := 0; i < len(lists)-1; i++ {
		if !(lists[i] < lists[i+1]) {
			panic("Error.")
		}
	}

}

func Test_GetgId1(t *testing.T) {
	gid1 := NewGIdImp()
	if gid1.GetId() != 0 {
		panic("Error.")
	}
	list1 := []uint64{}
	for i := uint64(0); i < 1000*1000; i++ {
		id1 := gid1.GenerateId()
		if gid1.GetId() != id1 {
			panic("Error.")
		}
		list1 = append(list1, gid1.GetId())
	}
	_ = list1
}
