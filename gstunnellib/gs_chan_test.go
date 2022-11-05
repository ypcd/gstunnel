package gstunnellib

import "testing"

func Test_chan_close_1(t *testing.T) {
	c1 := make(chan []byte)
	CloseChan(c1)
}

func noTest_chan_close_2(t *testing.T) {
	c1 := make(chan []byte)
	close(c1)
	CloseChan(c1)

	c2 := make(chan []byte, 10)
	c2 <- []byte{1}
	close(c2)
	CloseChan(c2)
}
