package gstunnellib

import (
	"bytes"
	"errors"
	"io"
	"net"
)

func IsTheVersionConsistent_sendEx(dst net.Conn, apack GsPack, wlent *int64, sendbuf *bytes.Buffer) error {
	if true {
		buf := apack.IsTheVersionConsistent()
		defer func() {
			if sendbuf != nil {
				io.Copy(sendbuf, bytes.NewReader(buf))
			}
		}()
		//tmr.Boot()
		//ChangeCryKey_Total += 1
		//outf2.Write(buf)
		for {
			if len(buf) > 0 {
				wlen, err := dst.Write(buf)
				*wlent = *wlent + int64(wlen)
				if wlen == 0 {
					return errors.New("wlen == 0")
				}
				if err != nil && wlen <= 0 {
					continue
				}
				if len(buf) == wlen {
					break
				}
				buf = buf[wlen:]
			} else {
				break
			}
		}
	}
	return nil
}

func IsTheVersionConsistent_send(dst net.Conn, apack GsPack, wlent *int64) error {
	return IsTheVersionConsistent_sendEx(dst, apack, wlent, nil)
}

func ChangeCryKey_sendEX(dst net.Conn, apack GsPack, ChangeCryKey_Total *int, wlent *int64, sendbuf *bytes.Buffer) error {

	buf := apack.ChangeCryKey()
	defer func() {
		if sendbuf != nil {
			io.Copy(sendbuf, bytes.NewReader(buf))
		}
	}()
	//tmr.Boot()
	*ChangeCryKey_Total += 1
	//outf2.Write(buf)
	for {
		if len(buf) > 0 {
			wlen, err := dst.Write(buf)
			*wlent = *wlent + int64(wlen)
			if wlen == 0 {
				return errors.New("wlen == 0")
			}
			if err != nil && wlen <= 0 {
				continue
			}
			if len(buf) == wlen {
				break
			}
			buf = buf[wlen:]
		} else {
			break
		}
	}
	return nil
}

func ChangeCryKey_send(dst net.Conn, apack GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return ChangeCryKey_sendEX(dst, apack, ChangeCryKey_Total, wlent, nil)
}
