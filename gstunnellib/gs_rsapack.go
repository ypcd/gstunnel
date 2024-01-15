package gstunnellib

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gspackoper"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type IGsRSAPack interface {
	IGSPack
	ClientPublicKeyPack() []byte
	ChangeCryKeyFromGSTServer() ([]byte, string)
	ChangeCryKeyFromGSTClient() ([]byte, string)
	PackPOHello() []byte

	IsExistsClientKey() bool
	GetClientRSAKey() *gsrsa.RSA
	SetClientRSAKey(*gsrsa.RSA)
	GetServerRSAKey() *gsrsa.RSA
	UnpackEx([]byte) (*UnpackReV, error)
}

type gsRSAPackImp struct {
	aespack
	serverkey *gsrsa.RSA
	clientkey *gsrsa.RSA
}

type UnpackReV struct {
	OperType  uint32
	ClientKey *gsrsa.RSA
	GenData   []byte
	Version   string
	Key       []byte
}

func (u *UnpackReV) IsPOVersion() bool {
	return gspackoper.IsPOVersionFromUint32(u.OperType)
}

func (u *UnpackReV) IsPOGen() bool {
	return gspackoper.IsPOGenFromUint32(u.OperType)
}

func (u *UnpackReV) IsClientPubKey() bool {
	return gspackoper.IsClientPubKeyFromUint32(u.OperType)
}

func (u *UnpackReV) IsChangeCryKey() bool {
	return gspackoper.IsChangeCryKeyFromUint32(u.OperType)
}

func (u *UnpackReV) IsChangeCryKeyFromGSTServer() bool {
	return gspackoper.IsChangeCryKeyFromGSTServerFromUint32(u.OperType)
}

func (u *UnpackReV) IsChangeCryKeyFromGSTClient() bool {
	return gspackoper.IsChangeCryKeyFromGSTClientFromUint32(u.OperType)
}

func (u *UnpackReV) IsPOHello() bool {
	return gspackoper.IsPOHello(u.OperType)
}

func (rp *gsRSAPackImp) ClientPublicKeyPack() []byte {
	ciphertext := rp.serverkey.EncryptWithPublicKeyAnyLen(rp.clientkey.PublicKeyToBytes())

	jdata := gspackoper.NewPackPOClientPubKey(ciphertext)

	crydata := rp.a.encrypter(jdata)

	IsOk_PackDataLen_panic(crydata)

	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata
}

func (rp *gsRSAPackImp) getClientPubKeyFromPackOper(po *gspackoper.PackOper) *gsrsa.RSA {
	if po.IsClientPubKey() {
		clientpubkeydata := po.GetClientPubKeyData()
		dedata := rp.serverkey.DecryptWithPrivateKeyAnyLen(clientpubkeydata)
		clientpubkey := gsrsa.PublicKeyFromBytes(dedata)
		return gsrsa.NewRSAObjPub(clientpubkey)
	}
	panic("PO type is not ClientPubKey")
}

func (rp *gsRSAPackImp) setClientPubKeyFromPackOper(po *gspackoper.PackOper) {

	if po.IsClientPubKey() {
		clientpubkeydata := po.GetClientPubKeyData()
		dedata := rp.serverkey.DecryptWithPrivateKeyAnyLen(clientpubkeydata)
		clientpubkey := gsrsa.PublicKeyFromBytes(dedata)
		if rp.clientkey == nil {
			rp.clientkey = gsrsa.NewRSAObjPub(clientpubkey)
		} else {
			rp.clientkey.SetPubkey(clientpubkey)
		}
		return
	}
	panic("error")
}

func (rp *gsRSAPackImp) Unpack(data []byte) ([]byte, error) {
	var jdata []byte
	if G_Deep_debug {
		G_logger.Println("aespack unpack data hash:", gshash.GetSha256Hex(data))

		jdata := rp.a.decrypter(data[2:])
		G_logger.Println("aespack unpack aes decrypt hash:", gshash.GetSha256Hex(jdata))
	} else {
		jdata = rp.a.decrypter(data[2:])
	}
	//var pdata []byte
	var err error

	po1 := gspackoper.UnPack_Oper(jdata)

	err = po1.IsOk()
	if err != nil {
		return nil, err
	}

	if po1.IsPOGen() {
		return po1.Data, err
	}

	if po1.IsChangeCryKeyFromGSTServer() {
		crykey := po1.GetEncryKey()
		key := rp.clientkey.DecryptWithPrivateKeyAnyLen(crykey)
		err = rp.setKey(key)
		G_logger.Println("gstunnel is ChangeCryKey.")
		//pdata = nil
		return nil, err
	}

	if po1.IsChangeCryKeyFromGSTClient() {
		crykey := po1.GetEncryKey()
		key := rp.serverkey.DecryptWithPrivateKeyAnyLen(crykey)
		err = rp.setKey(key)
		G_logger.Println("gstunnel is ChangeCryKey.")
		//pdata = nil
		return nil, err
	}

	if po1.IsPOVersion() {
		if string(po1.OperData) == G_Version {
			G_logger.Println("gstunnel POVersion is ok.")
			//pdata = nil
			//err = nil
			return nil, nil
		} else {
			//pdata = nil
			err = errors.New("G_Version is error. " + "error version: " + string(po1.OperData))
			panic(err)
		}
		//return nil, err
	}

	if po1.IsClientPubKey() {
		rp.setClientPubKeyFromPackOper(po1)
		return nil, nil
	}

	panic("error")
	//return pdata, err
}

func (rp *gsRSAPackImp) UnpackEx(data []byte) (*UnpackReV, error) {

	var jdata []byte
	if G_Deep_debug {
		G_logger.Println("aespack unpack data hash:", gshash.GetSha256Hex(data))

		jdata := rp.a.decrypter(data[2:])
		G_logger.Println("aespack unpack aes decrypt hash:", gshash.GetSha256Hex(jdata))
	} else {
		jdata = rp.a.decrypter(data[2:])
	}
	//var pdata []byte
	var err error

	//up1 := &UnpackReV{}

	po1 := gspackoper.UnPack_Oper(jdata)

	err = po1.IsOk()
	if err != nil {
		return nil, err
	}

	if po1.IsPOGen() {
		return &UnpackReV{
			OperType: po1.OperType,
			GenData:  po1.Data}, err
	}

	if po1.IsChangeCryKeyFromGSTServer() {
		crykey := po1.GetEncryKey()
		key := rp.clientkey.DecryptWithPrivateKeyAnyLen(crykey)
		//err = rp.setKey(key)
		G_logger.Println("gstunnel is ChangeCryKey.")
		//pdata = nil
		//return nil, err
		return &UnpackReV{
			OperType: po1.OperType,
			Key:      key}, err
	}

	if po1.IsChangeCryKeyFromGSTClient() {
		crykey := po1.GetEncryKey()
		key := rp.serverkey.DecryptWithPrivateKeyAnyLen(crykey)
		//err = rp.setKey(key)
		G_logger.Println("gstunnel is ChangeCryKey.")
		//pdata = nil
		return &UnpackReV{
			OperType: po1.OperType,
			Key:      key}, err
	}

	if po1.IsPOVersion() {
		return &UnpackReV{
			OperType: po1.OperType,
			Version:  string(po1.OperData)}, err
	}

	if po1.IsClientPubKey() {
		return &UnpackReV{
			OperType:  po1.OperType,
			ClientKey: rp.getClientPubKeyFromPackOper(po1)}, nil
	}

	if po1.IsPOHello() {
		return &UnpackReV{
			OperType: po1.OperType}, nil
	}

	panic("gstunnel pack type is error")
	//return pdata, err
}

func (rp *gsRSAPackImp) ChangeCryKey() []byte {
	panic("ChangeCryKey() is not extis")
}

func (rp *gsRSAPackImp) ChangeCryKeyFromGSTServer() (packData []byte, outkey string) {

	key := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	crykey := rp.clientkey.EncryptWithPublicKeyAnyLen(key)

	jdata := gspackoper.JsonPacking_OperChangeKeyRSAFromGSTServer(crykey)
	crydata := rp.a.encrypter(jdata)
	//edata := base64.StdEncoding.EncodeToString(crydata)

	rp.setKey(key)

	IsOk_PackDataLen_panic(crydata)

	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata, string(key)
}

func (rp *gsRSAPackImp) ChangeCryKeyFromGSTClient() (packData []byte, outkey string) {

	key := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	crykey := rp.serverkey.EncryptWithPublicKeyAnyLen(key)

	jdata := gspackoper.JsonPacking_OperChangeKeyRSAFromGSTClient(crykey)
	crydata := rp.a.encrypter(jdata)
	//edata := base64.StdEncoding.EncodeToString(crydata)

	rp.setKey(key)

	IsOk_PackDataLen_panic(crydata)

	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata, string(key)
}

func (rp *gsRSAPackImp) IsExistsClientKey() bool {
	if rp.clientkey == nil {
		return false
	} else {
		return true
	}
}

func (rp *gsRSAPackImp) GetClientRSAKey() *gsrsa.RSA { return rp.clientkey.NewRSA() }

func (rp *gsRSAPackImp) SetClientRSAKey(ckey *gsrsa.RSA) { rp.clientkey = ckey.NewRSA() }

func (rp *gsRSAPackImp) GetServerRSAKey() *gsrsa.RSA { return rp.serverkey.NewRSA() }

func (ap *gsRSAPackImp) PackPOHello() []byte {

	jdata := gspackoper.NewPackPOHello()
	crydata := ap.a.encrypter(jdata)
	//edata := base64.StdEncoding.EncodeToString(crydata)

	IsOk_PackDataLen_panic(crydata)

	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata
}

func newRSAPackImp(key string, serverRSAKey, ClientRSAKey *gsrsa.RSA) *gsRSAPackImp {
	if len([]byte(key)) != gsbase.G_AesKeyLen {
		checkError_exit(
			fmt.Errorf("error: The key is not %d bytes", gsbase.G_AesKeyLen))
	}

	if ClientRSAKey == nil {
		return &gsRSAPackImp{
			aespack:   aespack{createAes([]byte(key))},
			serverkey: serverRSAKey.NewRSA(),
		}
	} else {
		return &gsRSAPackImp{
			aespack:   aespack{createAes([]byte(key))},
			serverkey: serverRSAKey.NewRSA(),
			clientkey: ClientRSAKey.NewRSA(),
		}
	}
}

func NewRSAPackImpWithGSTClient(key string, serverRSAKey, ClientRSAKey *gsrsa.RSA) *gsRSAPackImp {
	if len([]byte(key)) != gsbase.G_AesKeyLen {
		checkError_exit(
			fmt.Errorf("error: The key is not %d bytes", gsbase.G_AesKeyLen))
	}

	return &gsRSAPackImp{
		aespack:   aespack{createAes([]byte(key))},
		serverkey: serverRSAKey.NewRSAPub(),
		clientkey: ClientRSAKey.NewRSA(),
	}
}

func NewRSAPackImpWithGSTServer(key string, serverRSAKey *gsrsa.RSA) *gsRSAPackImp {
	if len([]byte(key)) != gsbase.G_AesKeyLen {
		checkError_exit(
			fmt.Errorf("error: The key is not %d bytes", gsbase.G_AesKeyLen))
	}

	return &gsRSAPackImp{
		aespack:   aespack{createAes([]byte(key))},
		serverkey: serverRSAKey.NewRSA(),
	}
}

/*
func newRSAPack(key string, serverRSAKey *gsrsa.RSA, ClientRSAKey *gsrsa.RSA) *gsRSAPackImp {
	if serverRSAKey.Pri != nil {
		return newRSAPack__(key, serverRSAKey)
	}
}

func NewGSRSAPackWithGSTServer(key string, serverPri *rsa.PrivateKey) IGsRSAPack {
	return newRSAPackWithGSTServer(key, serverPri)
}

func NewGSRSAPackWithGSTClient(key string, serverPub *rsa.PublicKey) IGsRSAPack {
	return newRSAPackWithGSTClient(key, serverPub)
}


func NewGSRSAPackImp(key string, serverRSAKey, clientRSAKey *gsrsa.RSA) IGsRSAPack {
	return newRSAPackImp(key, serverRSAKey, clientRSAKey)
}
*/
