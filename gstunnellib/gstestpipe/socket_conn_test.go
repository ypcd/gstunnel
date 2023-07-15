package gstestpipe

import (
	"fmt"
	"testing"
)

func Test_GetRandAddr(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(GetRandAddr())
	}
}
