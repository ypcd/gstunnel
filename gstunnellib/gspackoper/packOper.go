package gspackoper

import (
	"bytes"
	"crypto/sha256"
	"errors"

	//"gstunnellib"

	//"gstunnellib"

	//"gstunnellib"
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"

	"google.golang.org/protobuf/proto"
)

const g_version string = gsbase.G_Version

var debug_tag bool

var p func(...interface{}) (int, error)

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

//type POtype uint32

const (
	POBegin                     uint32 = 1<<5 - 2
	POGenOper                   uint32 = 1<<5 - 1
	POChangeCryKey              uint32 = 1<<6 - 1
	POVersion                   uint32 = 1<<7 - 1
	POClientPubKey              uint32 = 1<<7 - 2
	POChangeCryKeyFromGSTServer uint32 = 1<<7 - 3
	POChangeCryKeyFromGSTClient uint32 = 1<<7 - 4
	POHello                     uint32 = 1<<7 - 5 //test use, Empty package
	POEnd                       uint32 = 1 << 7
)

type PackOper struct {
	packOperPro
}

func GetSha256_32bytes(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func (po *PackOper) GetSha256_old() []byte {
	v1 := GetSha256_32bytes(
		append(
			gsrand.Intu8ToBytes(uint8(po.OperType)), append(po.OperData, append(po.Data, po.Rand...)...)...,
		),
	)
	return v1[:]
}

func (po *PackOper) GetSha256() []byte {
	return po.GetSha256_nobuf()
}

// GetSha256_buf速度比GetSha256_old快15%-20%。
func (po *PackOper) GetSha256_buf() []byte {

	var buf bytes.Buffer
	buf.Write(gsrand.Intu8ToBytes(uint8(po.OperType)))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(po.Rand)

	sha := GetSha256_32bytes(buf.Bytes())

	return sha[:]
}

// 为了更好的安全性，默认情况下没有使用GetSha256_pool函数。
// GetSha256_pool的速度比GetSha256_buf的速度快10%-15%。
func (po *PackOper) GetSha256_pool() []byte {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.Write(gsrand.Intu8ToBytes(uint8(po.OperType)))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(po.Rand)

	h1 := GetSha256_32bytes(buf.Bytes())

	return h1[:]
}

// GetSha256_nobuf的速度比GetSha256_buf的速度快18%-20%。
func (po *PackOper) GetSha256_nobuf() []byte {
	hash := sha256.New()
	hash.Write(gsrand.Intu8ToBytes(uint8(po.OperType)))
	hash.Write(po.OperData)
	hash.Write(po.Data)
	hash.Write(po.Rand)

	return hash.Sum(nil)
}

func (po *PackOper) IsOk() error {
	if po.OperType <= POBegin || po.OperType >= POEnd {
		return errors.New("PackOper OperType is error")
	}

	h1 := po.GetSha256()

	if bytes.Equal(po.HashHex, h1) {
		return nil
	} else {
		return errors.New("the packoper hash is inconsistent")
	}
}

func (po *PackOper) IsPOVersion() bool {
	return po.OperType == POVersion
}

func (po *PackOper) IsPOGen() bool {
	return po.OperType == POGenOper
}

func (po *PackOper) IsClientPubKey() bool {
	return (po.OperType == POClientPubKey)
}

func (po *PackOper) IsChangeCryKey() bool {
	return po.OperType == POChangeCryKey
}

func (po *PackOper) IsChangeCryKeyFromGSTServer() bool {
	return po.OperType == POChangeCryKeyFromGSTServer
}

func (po *PackOper) IsChangeCryKeyFromGSTClient() bool {
	return po.OperType == POChangeCryKeyFromGSTClient
}

func (po *PackOper) IsPOHello() bool {
	return po.OperType == POHello
}

func (po *PackOper) GetEncryKey() []byte {
	if po.IsChangeCryKeyFromGSTServer() || po.IsChangeCryKeyFromGSTClient() || po.IsChangeCryKey() {
		return po.OperData
	}
	panic("GetEncryKey is error")
}

func (po *PackOper) GetClientPubKeyData() []byte {
	if po.IsClientPubKey() {
		return po.OperData
	}
	panic("GetClientPubKey is error")
}

func IsChangeCryKey(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POChangeCryKey)
}

func IsPOVersion(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POVersion)
}

func IsPOGen(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POGenOper)
}

func IsClientPubKey(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POClientPubKey)
}

func IsChangeCryKeyFromGSTServer(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POChangeCryKeyFromGSTServer)
}

func IsChangeCryKeyFromGSTClient(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POChangeCryKeyFromGSTClient)
}

func IsPOVersionFromUint32(OperType uint32) bool {
	return OperType == POVersion
}

func IsPOGenFromUint32(OperType uint32) bool {
	return OperType == POGenOper
}

func IsClientPubKeyFromUint32(OperType uint32) bool {
	return (OperType == POClientPubKey)
}

func IsChangeCryKeyFromUint32(OperType uint32) bool {
	return OperType == POChangeCryKey
}

func IsChangeCryKeyFromGSTServerFromUint32(OperType uint32) bool {
	return OperType == POChangeCryKeyFromGSTServer
}

func IsChangeCryKeyFromGSTClientFromUint32(OperType uint32) bool {
	return OperType == POChangeCryKeyFromGSTClient
}
func IsPOHello(OperType uint32) bool {
	return OperType == POHello
}

func getEncryKey(Data []byte) []byte {
	return (UnPack_Oper(Data).OperData)
}

func createPackOperGen(data []byte) *PackOper {

	pd := PackOper{packOperPro: packOperPro{
		OperType: POGenOper,
		Data:     data,
		Rand:     gsrand.GetRDBytes(int(128 + 256*gsrand.GetRDF64())),
	}}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func createPackOperChangeKey() (*PackOper, []byte) {

	key := gsrand.GetRDBytes(gsbase.G_AesKeyLen)

	pd := PackOper{packOperPro: packOperPro{
		OperType: POChangeCryKey,
		OperData: []byte(key),
		//Rand:     gsrand.GetRDBytes(8),
		Rand: gsrand.GetRDBytes(int(1024 + 512*gsrand.GetRDF64())),
	}}

	pd.HashHex = pd.GetSha256()
	return &pd, key
}

func createPackOperChangeKeyRSAFromGSTServer(crykey []byte) *PackOper {

	//key := gsrand.GetRDBytes(gsbase.G_AesKeyLen)

	pd := PackOper{packOperPro: packOperPro{
		OperType: POChangeCryKeyFromGSTServer,
		OperData: crykey,
		//Rand:     gsrand.GetRDBytes(8),
		Rand: gsrand.GetRDBytes(int(1024 + 512*gsrand.GetRDF64())),
	}}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func createPackOperChangeKeyRSAFromGSTClient(crykey []byte) *PackOper {

	//key := gsrand.GetRDBytes(gsbase.G_AesKeyLen)

	pd := PackOper{packOperPro: packOperPro{
		OperType: POChangeCryKeyFromGSTClient,
		OperData: crykey,
		//Rand:     gsrand.GetRDBytes(8),
		Rand: gsrand.GetRDBytes(int(1024 + 512*gsrand.GetRDF64())),
	}}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func JsonPacking_OperChangeKey() (data []byte, key []byte) {

	pd, rekey := createPackOperChangeKey()

	re, err := proto.Marshal(&pd.packOperPro)
	gserror.CheckError_panic(err)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re, rekey
}

func JsonPacking_OperChangeKeyRSAFromGSTServer(crykey []byte) []byte {

	pd := createPackOperChangeKeyRSAFromGSTServer(crykey)

	re, err := proto.Marshal(&pd.packOperPro)
	gserror.CheckError_panic(err)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

func JsonPacking_OperChangeKeyRSAFromGSTClient(crykey []byte) []byte {

	pd := createPackOperChangeKeyRSAFromGSTClient(crykey)

	re, err := proto.Marshal(&pd.packOperPro)
	gserror.CheckError_panic(err)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

//func JsonPacking_OperChangeKey_rsa(pub *rsa.PublicKey) []byte {}

func createPackOperVersion() *PackOper {

	//rd := gsrand.GetRDBytes(int(1024 + 512*gsrand.GetRDF64()))
	//rd := gsrand.GetRDBytes(8)

	pd := PackOper{packOperPro: packOperPro{
		OperType: POVersion,
		OperData: []byte(g_version),
		Rand:     gsrand.GetRDBytes(int(1024 + 512*gsrand.GetRDF64())),
	}}

	pd.HashHex = pd.GetSha256()

	return &pd
}

func createClientPubKey(keyCry []byte) *PackOper {
	pd := PackOper{packOperPro: packOperPro{
		OperType: POClientPubKey,
		OperData: keyCry,
		Rand:     gsrand.GetRDBytes(int(1024 + 512*gsrand.GetRDF64())),
	}}
	pd.HashHex = pd.GetSha256()
	return &pd
}

func NewPackPOVersion() []byte {

	pd := createPackOperVersion()

	re, err := proto.Marshal(&pd.packOperPro)
	gserror.CheckError_panic(err)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

func NewPackPOClientPubKey(keyCry []byte) []byte {

	pd := createClientPubKey(keyCry)

	re, err := proto.Marshal(&pd.packOperPro)
	gserror.CheckError_panic(err)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

func UnPack_Oper(data []byte) *PackOper {
	var msg PackOper

	err := proto.Unmarshal(data, &msg.packOperPro)
	if err != nil {
		panic(err)
	}
	//msg.packOperPro = msg2
	if debug_tag {
		p(msg)
		p("json: ", string(data))
	}

	return &msg
}

func jsonUnPack_OperGen(data []byte) ([]byte, error) {

	p1 := UnPack_Oper(data)

	return p1.Data, p1.IsOk()
}

func NewPackPOGen(data []byte) []byte {
	pd := createPackOperGen(data)

	re, err := proto.Marshal(&pd.packOperPro)
	gserror.CheckError_panic(err)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}
func createPackPOHello() *PackOper {

	pd := PackOper{packOperPro: packOperPro{
		OperType: POHello,
		Data:     nil,
		Rand:     gsrand.GetRDBytes(int(128 + 256*gsrand.GetRDF64())),
	}}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func NewPackPOHello() []byte {
	pd := createPackPOHello()

	re, err := proto.Marshal(&pd.packOperPro)
	gserror.CheckError_panic(err)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

func jsonUnpack(data []byte) ([]byte, error) {
	/*
		if begin_Dinfo == 0 {
			fmt.Println("Pack_Rand use.")
			begin_Dinfo = 1
		}
	*/
	return jsonUnPack_OperGen(data)
}
