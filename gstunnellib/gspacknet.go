package gstunnellib

import (
	"errors"
)

// net buffer
// encrypter decrypter
// pack unpack
type GsPackNet interface {
	WriteEncryData([]byte) error
	GetDecryData() ([]byte, error)

	Packing([]byte) []byte
	Unpack([]byte) ([]byte, error)

	ChangeCryKey() []byte
	IsTheVersionConsistent() []byte
}

type gsPackNetImp struct {
	apack GsPack
	buf   []byte
}

func NewGsPackNet(key string) GsPackNet {
	return &gsPackNetImp{apack: NewGsPack(key)}
}

func (pn *gsPackNetImp) WriteEncryData(data []byte) error {
	if len(data) <= 0 {
		return errors.New("len(data) <= 0.")
	}
	pn.buf = append(pn.buf, data...)
	return nil
}
func (pn *gsPackNetImp) GetDecryData() ([]byte, error) {
	rn, fdbl := Find0(pn.buf)
	var wbuf []byte
	if fdbl {
		wbuf = pn.buf[:rn]
		pn.buf = pn.buf[rn+1:]

		wbuf, err := pn.apack.Unpack(wbuf)
		checkError_panic(err)
		return wbuf, nil
	} else {
		return nil, nil
	}
}

func (pn *gsPackNetImp) Packing(data []byte) []byte {
	return pn.apack.Packing(data)
}
func (pn *gsPackNetImp) Unpack(data []byte) ([]byte, error) {
	return pn.apack.Unpack(data)
}

func (pn *gsPackNetImp) ChangeCryKey() []byte {
	return pn.apack.ChangeCryKey()
}
func (pn *gsPackNetImp) IsTheVersionConsistent() []byte {
	return pn.apack.IsTheVersionConsistent()
}
