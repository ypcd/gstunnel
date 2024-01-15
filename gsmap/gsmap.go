package gsmap

import (
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type gsmap struct {
	data sync.Map
}

func (m *gsmap) Add(key, value any) { m.data.Store(key, value) }
func (m *gsmap) Delete(key any)     { m.data.Delete(key) }

type GsMapStr_Str struct {
	gsmap
}

func NewGSMapStr_Str() *GsMapStr_Str {
	return &GsMapStr_Str{
		gsmap{data: sync.Map{}},
	}
}

func (m *GsMapStr_Str) Get(key string) string {
	v1, ok := m.data.Load(key)
	if !ok {
		panic("map load is error")
	}
	k2, ok := v1.(string)
	if !ok {
		panic("v1.(string) is error")
	}
	return k2
}

type GsMapStr_RSA struct {
	gsmap
}

func NewGSMapStr_RSA() *GsMapStr_RSA {
	return &GsMapStr_RSA{
		gsmap{data: sync.Map{}},
	}
}

func (m *GsMapStr_RSA) isOk_InputType(key, value any) {
	if key != nil {
		_, ok := key.(string)
		if !ok {
			panic("key.(T) is error")
		}
	}
	if value != nil {
		_, ok := value.(*gsrsa.RSA)
		if !ok {
			panic("value.(T) is error")
		}
	}
}

func (m *GsMapStr_RSA) Add(key, value any) {
	m.isOk_InputType(key, value)
	m.data.Store(key, value)
}
func (m *GsMapStr_RSA) Delete(key any) {
	m.isOk_InputType(key, nil)
	m.data.Delete(key)
}

func (m *GsMapStr_RSA) Get(key string) (*gsrsa.RSA, bool) {
	v1, ok := m.data.Load(key)
	if !ok {
		return nil, ok
	}
	if v1 == nil {
		return nil, ok
	}
	k2, ok := v1.(*gsrsa.RSA)
	if !ok {
		panic("v1.(T) is error")
	}
	return k2, true
}
