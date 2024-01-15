package gsobj

import (
	"bytes"
	"io"
	"net"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func VersionPack_sendEx(dst net.Conn, apack iGSTObj, wlent *int64, GetSendbuf *bytes.Buffer) error {

	if true {
		buf := apack.PackVersion()

		if GetSendbuf != nil {
			io.Copy(GetSendbuf, bytes.NewBuffer(buf))
		}

		wlen, err := gstunnellib.NetConnWriteAll(dst, buf)
		*wlent += int64(wlen)
		if gstunnellib.IsErrorNetUsually(err) {
			gserror.CheckError_panic(err)
			return err
		} else {
			gserror.CheckError_panic(err)
		}
	}
	return nil
}

func VersionPack_send(obj *GSTObj) error {
	return VersionPack_sendEx(obj.dst, obj, &obj.Wlent, nil)
}

/*
func changeCryKey_sendEX(dst net.Conn, apack IGSPack, ChangeCryKey_Total *int, wlent *int64, GetSendbuf *bytes.Buffer) error {

		buf := apack.ChangeCryKey()

		if GetSendbuf != nil {
			io.Copy(GetSendbuf, bytes.NewReader(buf))
		}
		//tmr.Boot()
		*ChangeCryKey_Total += 1
		//outf2.Write(buf)
		wlen, err := gstunnellib.NetConnWriteAll(dst, buf)
		*wlent += int64(wlen)
		if gstunnellib.IsErrorNetUsually(err) {
			gserror.CheckError_panic(err)
			return err
		} else {
			gserror.CheckError_panic(err)
		}
		return nil
	}
*/
func changeCryKey_sendEX_fromClient(dst net.Conn, apack iGSTObj, ChangeCryKey_Total *int, wlent *int64, GetSendbuf *bytes.Buffer) error {

	buf, _ := apack.ChangeCryKeyFromGSTClient()

	if GetSendbuf != nil {
		io.Copy(GetSendbuf, bytes.NewReader(buf))
	}
	//tmr.Boot()
	*ChangeCryKey_Total += 1
	//outf2.Write(buf)
	wlen, err := gstunnellib.NetConnWriteAll(dst, buf)
	*wlent += int64(wlen)
	if gstunnellib.IsErrorNetUsually(err) {
		gserror.CheckError_panic(err)
		return err
	} else {
		gserror.CheckError_panic(err)
	}
	return nil
}

func changeCryKey_sendEX_fromServer(dst net.Conn, apack iGSTObj, ChangeCryKey_Total *int, wlent *int64, GetSendbuf *bytes.Buffer) error {

	buf, _ := apack.ChangeCryKeyFromGSTServer()

	if GetSendbuf != nil {
		io.Copy(GetSendbuf, bytes.NewReader(buf))
	}
	//tmr.Boot()
	*ChangeCryKey_Total += 1
	//outf2.Write(buf)
	wlen, err := gstunnellib.NetConnWriteAll(dst, buf)
	*wlent += int64(wlen)
	if gstunnellib.IsErrorNetUsually(err) {
		gserror.CheckError_panic(err)
		return err
	} else {
		gserror.CheckError_panic(err)
	}
	return nil
}

/*
	func ChangeCryKey_send(dst net.Conn, apack IGSPack, ChangeCryKey_Total *int, wlent *int64) error {
		return changeCryKey_sendEX(dst, apack, ChangeCryKey_Total, wlent, nil)
	}
*/
func changeCryKey_send_fromClient(obj *GSTObj) error {
	return changeCryKey_sendEX_fromClient(obj.dst, obj, &obj.ChangeCryKey_Total, &obj.Wlent, nil)
}

func changeCryKey_send_fromServer(obj *GSTObj) error {
	return changeCryKey_sendEX_fromServer(obj.dst, obj, &obj.ChangeCryKey_Total, &obj.Wlent, nil)
}

func newGSTObj_net_test(dst net.Conn) *GSTObj {
	return &GSTObj{
		dst:   dst,
		apack: gstunnellib.NewGSRSAPackNetImp(gsbase.G_AesKeyDefault, gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen), gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)),
	}
}
