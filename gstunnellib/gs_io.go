package gstunnellib

import (
	"net"
	"os"
)

func NetConnWriteAll(dst net.Conn, buf []byte) (int64, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	var wlen int = 0
	for {
		wsz, err := dst.Write(buf)
		wlen += wsz
		if err != nil {
			return 0, err
		}
		if wsz == len(buf) {
			return int64(wlen), err
		} else if wsz > len(buf) {
			panic("error wlen>len(buf)")
		}
		buf = buf[wsz:]
	}
}

// 比io.copy快12.9%
func netConnWriteAll_test(dst *os.File, buf []byte) (int64, error) {
	var wlen int = 0
	//buflen := len(buf)
	for {
		wsz, err := dst.Write(buf)
		wlen += wsz
		if err != nil {
			return 0, err
		}
		if wsz == len(buf) {
			return int64(wlen), err
		} else if wsz > len(buf) {
			panic("error wlen>len(buf)")
		}
		buf = buf[wsz:]
	}
}
