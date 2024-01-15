package gstunnellib

import (
	"bytes"
	"log"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_Aest(t *testing.T) {

	fbuf := gsrand.GetRDBytes(1024 * 1024)

	a1 := createAes(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	tmp := a1.encrypter(fbuf)
	outbuf := a1.decrypter(tmp)

	if bytes.Equal(fbuf, outbuf) {

	} else {
		t.Error()
	}

}

func Test_Aest2(t *testing.T) {

	a1 := createAes(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	for i := 0; i < 1000; i++ {
		fbuf := gsrand.GetRDBytes(1024 * 1024)

		tmp := a1.encrypter(fbuf)
		outbuf := a1.decrypter(tmp)

		if bytes.Equal(fbuf, outbuf) {

		} else {
			t.Error()
		}
	}
}

func Test_Aest3(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fbuf := gsrand.GetRDBytes(1024 * 1024)

		a1 := createAes(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

		tmp := a1.encrypter(fbuf)
		outbuf := a1.decrypter(tmp)

		if bytes.Equal(fbuf, outbuf) {

		} else {
			t.Error()
		}
	}
}

func aest() {

	fbuf := gsrand.GetRDBytes(1024)

	a1 := createAes(gsrand.GetRDBytes(gsbase.G_AesKeyLen))
	tmp := a1.encrypter(fbuf)
	outbuf := a1.decrypter(tmp)

	if bytes.Equal(fbuf, outbuf) {
		p("ok.", getrand())

	} else {
		log.Fatal("Error")
	}
	wlist.Done()
}

func Test_MtAest(t *testing.T) {
	mtF(aest)
}
