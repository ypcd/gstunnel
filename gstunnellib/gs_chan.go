package gstunnellib

import "errors"

func CloseChan(c interface{}) {
	chan1, ok := c.(chan []byte)
	if !ok {
		checkError_panic(errors.New("Error: The c is not chan []byte."))
	}
	defer Panic_Recover(logger)
	select {
	case _, ok = <-chan1:
		if ok {
			close(chan1)
		}
	default:
		close(chan1)
	}
}
