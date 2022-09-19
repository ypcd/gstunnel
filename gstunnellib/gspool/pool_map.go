package gspool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type pool_map struct {
	data       map[int64][]byte
	data_mutex sync.Mutex
	bytesSize  int
	gid        int64
	maxsize    int

	count_newBytes      int
	count_put_not_add   int
	count_get_from_data int
	count_get           int
}

func NewPoolMap(inbytesSize, inmaxsize int) *pool_map {
	pl := pool_map{
		data:      make(map[int64][]byte, 0),
		bytesSize: inbytesSize,
		gid:       0,
		maxsize:   inmaxsize,
	}
	return &pl
}

func (pl *pool_map) newBytes() []byte {
	pl.count_newBytes++
	return make([]byte, pl.bytesSize)
}
func (pl *pool_map) getNewId() int64 {
	return atomic.AddInt64(&pl.gid, 1)
}
func (pl *pool_map) Get() []byte {
	pl.data_mutex.Lock()
	defer pl.data_mutex.Unlock()
	pl.count_get++

	if len(pl.data) == 0 {
		return pl.newBytes()
	}
	pl.count_get_from_data++

	var k int64
	var value []byte
	for k, value = range pl.data {
		break
	}
	delete(pl.data, k)
	return value
}
func (pl *pool_map) Put(indata []byte) {
	pl.data_mutex.Lock()
	defer pl.data_mutex.Unlock()
	if len(pl.data) < pl.maxsize {
		pl.data[pl.getNewId()] = indata
		return
	} else {
		pl.count_put_not_add++
		return
	}
}
func (pl *pool_map) Size() int {
	return len(pl.data)
}
func (pl *pool_map) print() {
	fmt.Printf("%+v\n", pl)
}

func (pl *pool_map) ClearAll() {
	pl.data = make(map[int64][]byte, 0)
}
