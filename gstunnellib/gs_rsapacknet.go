package gstunnellib

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type IGSRSAPackNet interface {
	IGSPackNet

	ClientPublicKeyPack() []byte
	ChangeCryKeyFromGSTServer() ([]byte, string)
	ChangeCryKeyFromGSTClient() ([]byte, string)
	PackPOHello() []byte

	IsExistsClientKey() bool
	GetClientRSAKey() *gsrsa.RSA
	SetClientRSAKey(*gsrsa.RSA)
	GetServerRSAKey() *gsrsa.RSA
	UnpackEx([]byte) (*UnpackReV, error)
	UnpackOneGSTPackFromNetConn(src net.Conn, timeout time.Duration) (*UnpackReV, error)
}

type gsRSAPackNetImp struct {
	gsPackNetImp
}

func (rp *gsRSAPackNetImp) ClientPublicKeyPack() []byte {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.ClientPublicKeyPack()
}

func (rp *gsRSAPackNetImp) ChangeCryKeyFromGSTServer() (packData []byte, outkey string) {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.ChangeCryKeyFromGSTServer()
}

func (rp *gsRSAPackNetImp) ChangeCryKeyFromGSTClient() (packData []byte, outkey string) {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.ChangeCryKeyFromGSTClient()
}

func (rp *gsRSAPackNetImp) PackPOHello() []byte {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.PackPOHello()
}

func (rp *gsRSAPackNetImp) IsExistsClientKey() bool {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.IsExistsClientKey()
}

func (rp *gsRSAPackNetImp) GetClientRSAKey() *gsrsa.RSA {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.GetClientRSAKey()
}

func (rp *gsRSAPackNetImp) SetClientRSAKey(ckey *gsrsa.RSA) {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	it.SetClientRSAKey(ckey)
}

func (rp *gsRSAPackNetImp) GetServerRSAKey() *gsrsa.RSA {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.GetServerRSAKey()
}

func (rp *gsRSAPackNetImp) UnpackEx(data []byte) (*UnpackReV, error) {
	it, ok := rp.apack.(IGsRSAPack)
	if !ok {
		panic("error")
	}
	return it.UnpackEx(data)
}

func (rp *gsRSAPackNetImp) readOneGSTPackFromNetConn(src net.Conn, timeout time.Duration) ([]byte, error) {
	//panic("")
	sizedata := make([]byte, 2)
	src.SetReadDeadline(time.Now().Add(timeout))
	_, err := io.ReadAtLeast(src, sizedata, 2)
	if gserror.IsErrorNetUsually(err) {
		return nil, errors.New(fmt.Sprintf("(*gsRSAPackNetImp) readOneGSTPackFromNetConn()::io.ReadAtLeast():[%s]: %s\n", src.RemoteAddr().String(), err.Error()))
	} else {
		checkError_panic(err)
	}
	sz := rp.GetGSPackSize(sizedata)

	readbody := make([]byte, sz)
	src.SetReadDeadline(time.Now().Add(timeout))
	_, err = io.ReadAtLeast(src, readbody, len(readbody))
	if gserror.IsErrorNetUsually(err) {
		return nil, errors.New(fmt.Sprintf("(*gsRSAPackNetImp) readOneGSTPackFromNetConn()::io.ReadAtLeast():[%s]: %s\n", src.RemoteAddr().String(), err.Error()))
	} else {
		checkError_panic(err)
	}
	return append(sizedata, readbody...), nil
}

func (rp *gsRSAPackNetImp) UnpackOneGSTPackFromNetConn(src net.Conn, timeout time.Duration) (*UnpackReV, error) {
	redata, err := rp.readOneGSTPackFromNetConn(src, timeout)
	if err != nil {
		return nil, err
	}
	return rp.UnpackEx(redata)
}

/*
	func NewGSRSAPackNetImpWithGSTServer(key string, serverPri *rsa.PrivateKey) IGSRSAPackNet {
		return &gsRSAPackNetImp{
			gsPackNetImp{
				apack: NewGSRSAPack(key, serverPri),
			},
		}
	}
*/

func NewGSRSAPackNetImp(key string, serverRSAKey, clientRSAKey *gsrsa.RSA) *gsRSAPackNetImp {
	if clientRSAKey != nil {
		return &gsRSAPackNetImp{
			gsPackNetImp{
				apack: newRSAPackImp(key, serverRSAKey, clientRSAKey),
			},
		}
	} else {
		return &gsRSAPackNetImp{
			gsPackNetImp{
				apack: newRSAPackImp(key, serverRSAKey, nil),
			},
		}
	}
}

func NewGSRSAPackNetImpWithGSTClient(key string, serverRSAKey, clientRSAKey *gsrsa.RSA) IGSRSAPackNet {
	return &gsRSAPackNetImp{
		gsPackNetImp{
			apack: NewRSAPackImpWithGSTClient(key, serverRSAKey, clientRSAKey),
		},
	}
}

func NewGSRSAPackNetImpWithGSTServer(key string, serverRSAKey *gsrsa.RSA) IGSRSAPackNet {
	return &gsRSAPackNetImp{
		gsPackNetImp{
			apack: NewRSAPackImpWithGSTServer(key, serverRSAKey),
		},
	}
}
