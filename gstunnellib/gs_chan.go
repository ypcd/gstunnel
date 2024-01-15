package gstunnellib

import (
	"errors"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

func ChanClose(c interface{}) {
	chan1, ok := c.(chan []byte)
	if !ok {
		checkError_panic(errors.New("error: the c is not chan []byte"))
	}
	defer gserror.Panic_Recover(G_logger)
	select {
	case _, ok = <-chan1:
		if ok {
			close(chan1)
		}
	default:
		close(chan1)
	}
}

func ChanClean(c chan []byte) {
	for range c {
	}
}
