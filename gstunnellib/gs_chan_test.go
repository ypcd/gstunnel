package gstunnellib

import "testing"

func Test_chan1(t *testing.T) {
	c1 := make(chan []byte)
	CloseChan(c1)
}
