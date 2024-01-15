package gsobj

import "github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"

type iGSTObj interface {
	Close()

	//SetClientRSAKey(ckey *gsrsa.RSA)
	WriteEncryData(data []byte) error

	GetDecryData() ([]byte, error)

	ClientPublicKeyPack() []byte

	IsExistsClientKey() bool

	GetClientRSAKey() *gsrsa.RSA

	PackVersion() []byte

	Packing(data []byte) []byte

	VersionPack_send() error

	ChangeCryKey_send() error

	ChangeCryKeyFromGSTServer() (packData []byte, outkey string)

	ChangeCryKeyFromGSTClient() (packData []byte, outkey string)
}
