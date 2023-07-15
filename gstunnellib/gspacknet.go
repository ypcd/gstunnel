package gstunnellib

import (
	"encoding/binary"
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

	GetGSPackSize(data []byte) uint16
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
		return errors.New("len(data) <= 0")
	}
	pn.buf = append(pn.buf, data...)
	return nil
}

func (pn *gsPackNetImp) GetDecryData() ([]byte, error) {

	if len(pn.buf) < 2 {
		return nil, nil
	}
	var rebuf []byte
	var wbuf []byte
	for {

		rn := GetGSPackSize(pn.buf[:2])
		if int64(2+rn) > int64(len(pn.buf)) {
			return rebuf, nil
		}
		if int64(2+rn) == int64(len(pn.buf)) {
			wbuf = pn.buf[:2+rn]
			pn.buf = nil
		} else {
			wbuf = pn.buf[:2+rn]
			pn.buf = pn.buf[2+rn:]
		}
		wbuf, err := pn.apack.Unpack(wbuf)
		CheckError_panic(err)
		if len(wbuf) > 0 {
			rebuf = append(rebuf, wbuf...)
		}
		if len(pn.buf) < 2 {
			return rebuf, nil
		}
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

func (pn *gsPackNetImp) GetGSPackSize(data []byte) uint16 {
	return GetGSPackSize(data)
}

func GetGSPackSize(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}
