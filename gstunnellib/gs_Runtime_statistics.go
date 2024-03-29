package gstunnellib

import (
	"encoding/json"
	"errors"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

type IRuntime_statistics interface {
	//	AddTotalNetworkData(int)
	AddSrcTotalNetData_recv(int)
	AddSrcTotalNetData_send(int)
	AddServerTotalNetData_recv(int)
	AddServerTotalNetData_send(int)

	GetJson() ([]byte, error)
}

type net_data struct {
	NetData      uint64
	NetData_recv uint64
	NetData_send uint64
}

func (nd *net_data) AddNetData_recv(inv int) {
	if inv < 0 {
		return
	}
	atomic.AddUint64(&nd.NetData_recv, uint64(inv))
	atomic.AddUint64(&nd.NetData, uint64(inv))
}
func (nd *net_data) AddNetData_send(inv int) {
	if inv < 0 {
		return
	}
	atomic.AddUint64(&nd.NetData_send, uint64(inv))
	atomic.AddUint64(&nd.NetData, uint64(inv))
}

type runtime_statistics_data struct {
	Goroutines   int
	TotalNetData uint64

	Src    net_data
	Server net_data

	//us 1000us=1ms
	PauseTotalNs float64
	NumGC        uint32
}

func copy_atomic_netdata(src *net_data) net_data {
	return net_data{
		NetData:      atomic.LoadUint64(&src.NetData),
		NetData_recv: atomic.LoadUint64(&src.NetData_recv),
		NetData_send: atomic.LoadUint64(&src.NetData_send),
	}
}

func deepcopy_atomic_rsd(insrc *runtime_statistics_data) runtime_statistics_data {
	inGoroutines := int64(insrc.Goroutines)

	return runtime_statistics_data{
		Goroutines:   int(atomic.LoadInt64(&inGoroutines)),
		TotalNetData: atomic.LoadUint64(&insrc.TotalNetData),

		Src:    copy_atomic_netdata(&insrc.Src),
		Server: copy_atomic_netdata(&insrc.Server),
	}
}

type runtime_statistics_imp struct {
	runtime_statistics_data

	lock sync.Mutex
}

func NewRuntimeStatistics() IRuntime_statistics {
	return &runtime_statistics_imp{}
}
func (rs *runtime_statistics_imp) AddSrcTotalNetData_recv(inv int) {

	if inv < 0 {
		return
	}
	rs.Src.AddNetData_recv(inv)
	atomic.AddUint64(&rs.TotalNetData, uint64(inv))
}
func (rs *runtime_statistics_imp) AddSrcTotalNetData_send(inv int) {

	if inv < 0 {
		return
	}
	rs.Src.AddNetData_send(inv)
	atomic.AddUint64(&rs.TotalNetData, uint64(inv))
}
func (rs *runtime_statistics_imp) AddServerTotalNetData_recv(inv int) {

	if inv < 0 {
		return
	}
	rs.Server.AddNetData_recv(inv)
	atomic.AddUint64(&rs.TotalNetData, uint64(inv))
}
func (rs *runtime_statistics_imp) AddServerTotalNetData_send(inv int) {

	if inv < 0 {
		return
	}
	rs.Server.AddNetData_send(inv)
	atomic.AddUint64(&rs.TotalNetData, uint64(inv))
}

func (rs *runtime_statistics_imp) setNumGoroutine() {

	rs.Goroutines = runtime.NumGoroutine()
}

func (rs *runtime_statistics_imp) GetJson() ([]byte, error) {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.setNumGoroutine()

	mems := runtime.MemStats{}
	runtime.ReadMemStats(&mems)
	rs.PauseTotalNs = float64(mems.PauseTotalNs) / 1000.0
	rs.NumGC = mems.NumGC

	//data1 := rs.runtime_statistics_data
	//json.Marshal(data1)

	data := deepcopy_atomic_rsd(&rs.runtime_statistics_data)
	/*
		p1 := &rs.runtime_statistics_data
		p2 := &data
		_, _ = p1, p2
		s1 := fmt.Sprintf("%p %p", p1, p2)
		_ = s1
	*/
	re, err := json.Marshal(data)
	return re, err
}

type runtime_statistics_null struct{}

func NewRuntimeStatisticsNULL() IRuntime_statistics                    { return &runtime_statistics_null{} }
func (rs *runtime_statistics_null) AddSrcTotalNetData_recv(inv int)    {}
func (rs *runtime_statistics_null) AddSrcTotalNetData_send(inv int)    {}
func (rs *runtime_statistics_null) AddServerTotalNetData_recv(inv int) {}
func (rs *runtime_statistics_null) AddServerTotalNetData_send(inv int) {}
func (rs *runtime_statistics_null) GetJson() ([]byte, error)           { return nil, nil }

func RunGRuntimeStatistics_print(inlog *log.Logger, inruntstats IRuntime_statistics) {
	for {
		time.Sleep(time.Second * 10)
		re, err := inruntstats.GetJson()
		if err != nil {
			inlog.Println("Error:", err.Error())
		} else {
			inlog.Println(string(re))
		}
		runtstats_data, ok := inruntstats.(*runtime_statistics_imp)
		if !ok {
			gserror.CheckErrorEx(errors.New("error"), inlog)
		}

		inlog.Println("memstats.PauseTotalNs:", runtstats_data.PauseTotalNs, "us",
			"\n", "memstats.NumGC:", runtstats_data.NumGC)
	}
}
