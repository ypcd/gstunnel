package gstunnellib

import (
	"fmt"
	"testing"
	"time"
)

func Test_nettime1(t *testing.T) {
	nt := NewNetTimeImp()
	nt.Add(10)
	fmt.Println(nt.PrintString(), nt)

	ntn := NewNetTimeImpName("test")
	ntn.Add(10)
	fmt.Println(ntn.PrintString(), nt)
}

func Test_nettime2(t *testing.T) {
	nt := NewNetTimeImp()

	for i := 1; i < 101; i++ {
		nt.Add(time.Millisecond * time.Duration(i))
	}
	fmt.Println(nt.PrintString())

	ntimp, _ := nt.(*netTimeImp)
	if ntimp.sum != (time.Millisecond * 101 * 50) {
		panic("Error.")
	}
	if ntimp.count != 100 {
		panic("Error.")
	}
	if ntimp.min != 1*time.Millisecond {
		panic("Error.")
	}
	if ntimp.max != 100*time.Millisecond {
		panic("Error.")
	}
}
