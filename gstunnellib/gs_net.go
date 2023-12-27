package gstunnellib

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
)

/*
	func init() {
		fmt.Println("gs_net init().")
	}
*/
func isTheVersionConsistent_sendEx(dst net.Conn, apack GsPack, wlent *int64, GetSendbuf *bytes.Buffer) error {
	if true {
		buf := apack.IsTheVersionConsistent()

		if GetSendbuf != nil {
			io.Copy(GetSendbuf, bytes.NewBuffer(buf))
		}

		wlen, err := NetConnWriteAll(dst, buf)
		*wlent += int64(wlen)
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			CheckError_panic(err)
			return err
		} else {
			CheckError_panic(err)
		}
	}
	return nil
}

func IsTheVersionConsistent_send(dst net.Conn, apack GsPack, wlent *int64) error {
	return isTheVersionConsistent_sendEx(dst, apack, wlent, nil)
}

func changeCryKey_sendEX(dst net.Conn, apack GsPack, ChangeCryKey_Total *int, wlent *int64, GetSendbuf *bytes.Buffer) error {

	buf := apack.ChangeCryKey()

	if GetSendbuf != nil {
		io.Copy(GetSendbuf, bytes.NewReader(buf))
	}
	//tmr.Boot()
	*ChangeCryKey_Total += 1
	//outf2.Write(buf)
	wlen, err := NetConnWriteAll(dst, buf)
	*wlent += int64(wlen)
	if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
		CheckError_panic(err)
		return err
	} else {
		CheckError_panic(err)
	}
	return nil
}

func ChangeCryKey_send(dst net.Conn, apack GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return changeCryKey_sendEX(dst, apack, ChangeCryKey_Total, wlent, nil)
}
