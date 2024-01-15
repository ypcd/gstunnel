package gsmap

import (
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
)

type IGSMapLock interface {
	Add(key uint64, value *gsobj.GSTConn)
}

type gsmapLock struct {
	data      map[uint64]*gsobj.GSTConn
	lock_data sync.Mutex
}

func NewGSMapLock() *gsmapLock {
	return &gsmapLock{data: make(map[uint64]*gsobj.GSTConn)}
}

func (m *gsmapLock) Add(key uint64, value *gsobj.GSTConn) {
	m.lock_data.Lock()
	defer m.lock_data.Unlock()
	_, ok := m.data[key]
	if ok {
		panic("The key of map already exists")
	}
	m.data[key] = value
}
