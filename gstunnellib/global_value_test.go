package gstunnellib

import (
	"fmt"
	"sync"
	"testing"
)

func Test_global_value1(t *testing.T) {
	g1 := NewGlobalValuesImp()
	fmt.Println(g1.GetDebug())
	g1.SetDebug(true)
	fmt.Println(g1.GetDebug())
}

func Test_global_value_go(t *testing.T) {
	g1 := NewGlobalValuesImp()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			g1.GetDebug()
		}
	}()

	for i := 0; i < 100000; i++ {
		g1.SetDebug(true)
	}
	wg.Wait()
}
