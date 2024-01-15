package gstunnellib

import (
	"encoding/binary"
	"errors"
)

// net buffer
// encrypter decrypter
// pack unpack
type IGSPackNet interface {
	IGSPack

	WriteEncryData([]byte) error
	GetDecryData() ([]byte, error)
	GetDecryDataFormBytes([]byte) ([]byte, error)

	GetGSPackSize(data []byte) uint16
}

type gsPackNetImp struct {
	apack IGSPack
	buf   []byte
}

func NewGsPackNet(key string) IGSPackNet {
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
		checkError_panic(err)
		if len(wbuf) > 0 {
			rebuf = append(rebuf, wbuf...)
		}
		if len(pn.buf) < 2 {
			return rebuf, nil
		}
	}
}

func (pn *gsPackNetImp) GetDecryDataFormBytes(data []byte) ([]byte, error) {
	error := pn.WriteEncryData(data)
	checkErrorEx_panic(error)
	return pn.GetDecryData()
}

func (pn *gsPackNetImp) Packing(data []byte) []byte {
	return pn.apack.Packing(data)
}

func (pn *gsPackNetImp) Unpack(data []byte) ([]byte, error) {
	panic("unpack is not exist")
}

func (pn *gsPackNetImp) ChangeCryKey() []byte {
	return pn.apack.ChangeCryKey()
}
func (pn *gsPackNetImp) PackVersion() []byte {
	return pn.apack.PackVersion()
}

func (pn *gsPackNetImp) GetGSPackSize(data []byte) uint16 {
	return GetGSPackSize(data)
}

func GetGSPackSize(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}
