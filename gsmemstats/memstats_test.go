package gsmemstats

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_LookHeapMemAndGoRun1(t *testing.T) {
	look := NewLookHeapMemAndGoRun(log.Default())
	_ = look
	time.Sleep(time.Second * 4)
}

func Test_LookHeapMemAndGoRun2(t *testing.T) {
	look := NewLookHeapMemAndGoRun_nolog()
	_ = look
	time.Sleep(time.Second * 4)
}

func Test_LookHeapMemAndMaxMem1(t *testing.T) {
	look := NewLookHeapMemAndMaxMem()
	//time.Sleep(time.Millisecond * 150)
	fmt.Println(look.GetMaxMemInfo())
	fmt.Println(look.GetMaxHeapMemStatsMiB())
	v1 := make([]byte, 1024*1024*100)
	for i := range v1 {
		v1[i] = byte(i)
	}
	_ = v1
	time.Sleep(time.Millisecond * 150)
	fmt.Println(look.GetMaxMemInfo())
	fmt.Println(look.GetMaxHeapMemStatsMiB())
	fmt.Println(look.GetHeapMemStatsMiB())
}
