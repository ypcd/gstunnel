package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type LookHeapMem struct {
	memStats *runtime.MemStats
}

func NewLookHeapMem() *LookHeapMem {
	return &LookHeapMem{&runtime.MemStats{}}
}

func (l *LookHeapMem) GetHeapMemStats() string {
	runtime.ReadMemStats(l.memStats)
	return fmt.Sprintf("HeapSys: %d, HeapIdle: %d, HeapReleased: %d, HeapInuse: %d, HeapAlloc: %d, HeapObjects: %d",
		l.memStats.HeapSys, l.memStats.HeapIdle, l.memStats.HeapReleased, l.memStats.HeapInuse, l.memStats.HeapAlloc, l.memStats.HeapObjects)
}

func (l *LookHeapMem) GetHeapMemStatsMiB() string {
	runtime.ReadMemStats(l.memStats)
	return fmt.Sprintf("HeapSys: %dMiB, HeapIdle: %dMiB, HeapReleased: %dMiB, HeapInuse: %dMiB, HeapAlloc: %dMiB, HeapObjects: %d",
		l.memStats.HeapSys/1024/1024, l.memStats.HeapIdle/1024/1024, l.memStats.HeapReleased/1024/1024, l.memStats.HeapInuse/1024/1024, l.memStats.HeapAlloc/1024/1024, l.memStats.HeapObjects)
}

// GC stop-the-world
func (l *LookHeapMem) GetGCInfo() string {
	runtime.ReadMemStats(l.memStats)
	return fmt.Sprintf("GC: GCCPUFraction: %f, PauseTotalNs: %d  %dms, PauseTotalNs: %d, NumGC: %d", l.memStats.GCCPUFraction, l.memStats.PauseTotalNs, l.memStats.PauseTotalNs/1000/1000, l.memStats.PauseNs, l.memStats.NumGC)
}

func LookHeapMemFunc() {
	//timer_displayer := timerm.CreateTimer(time.Second * 3)

	memStats := NewLookHeapMem()
	defer logger_test.Println(memStats.GetGCInfo())

	time.Sleep(time.Millisecond * 100)
	//logger_test.Println("MemStats:", memStats)
	//	logger_test.Printf("HeapSys: %d, HeapIdle: %d, HeapReleased: %d, HeapInuse: %d, HeapAlloc: %d, HeapObjects: %d\n",
	//		memStats.HeapSys, memStats.HeapIdle, memStats.HeapReleased, memStats.HeapInuse, memStats.HeapAlloc, memStats.HeapObjects)
	for {

		//logger_test.Println("MemStats:", memStats)
		//		logger_test.Printf("HeapSys: %d, HeapIdle: %d, HeapReleased: %d, HeapInuse: %d, HeapAlloc: %d, HeapObjects: %d\n",
		//			memStats.HeapSys, memStats.HeapIdle, memStats.HeapReleased, memStats.HeapInuse, memStats.HeapAlloc, memStats.HeapObjects)
		logger_test.Println(memStats.GetHeapMemStatsMiB())
		logger_test.Println(memStats.GetGCInfo())
		time.Sleep(time.Second * 3)
	}
}

func init() {
	//go LookHeapMemFunc()
}

func Test_server_socket(t *testing.T) {
	//debug.SetGCPercent(50)
	inTest_server_socket(t, false, 200)
}

func Test_server_socket_mtg(t *testing.T) {
	inTest_server_socket(t, true, 200)
}

func Test_server_socket_mt(t *testing.T) {
	inTest_server_socket_mt(t, false, 100)
}

func Test_server_socket_mt_mtg(t *testing.T) {
	inTest_server_socket_mt(t, true, 100)
}

func noUseTest_server_socket_mt_old_listen(t *testing.T) {
	inTest_server_socket_mt_old_listen(t, false, 100)
}
