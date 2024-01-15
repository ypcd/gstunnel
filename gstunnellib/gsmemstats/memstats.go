package gsmemstats

import (
	"fmt"
	"log"
	"runtime"
	"sync"
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
	return fmt.Sprintf("GC: GCCPUFraction: %f, PauseTotalNs: %d  %dms, NumGC: %d, PauseNs: %d,", l.memStats.GCCPUFraction, l.memStats.PauseTotalNs, l.memStats.PauseTotalNs/1000/1000, l.memStats.NumGC, l.memStats.PauseNs)
}

type LookHeapMemAndGoPrintInfo struct {
	LookHeapMem
	inlogger   *log.Logger
	running_go bool
}

func NewLookHeapMemAndGoRun(inlog *log.Logger) *LookHeapMemAndGoPrintInfo {
	look := &LookHeapMemAndGoPrintInfo{LookHeapMem: LookHeapMem{&runtime.MemStats{}},
		inlogger: inlog,
	}
	look.run()
	return look
}

func NewLookHeapMemAndGoRun_nolog() *LookHeapMemAndGoPrintInfo {
	look := &LookHeapMemAndGoPrintInfo{LookHeapMem: LookHeapMem{&runtime.MemStats{}}}
	look.run()
	return look
}

func LookHeapMemFunc(look *LookHeapMemAndGoPrintInfo) {
	memStats := &look.LookHeapMem
	defer fmt.Println(memStats.GetGCInfo())

	time.Sleep(time.Millisecond * 100)
	for {
		fmt.Println(memStats.GetHeapMemStatsMiB())
		fmt.Println(memStats.GetGCInfo())
		time.Sleep(time.Second * 3)
	}
}

func LookHeapMemFunc_log(look *LookHeapMemAndGoPrintInfo) {
	logger := look.inlogger
	memStats := &look.LookHeapMem
	defer logger.Println(memStats.GetGCInfo())

	time.Sleep(time.Millisecond * 100)
	for {
		logger.Println(memStats.GetHeapMemStatsMiB())
		logger.Println(memStats.GetGCInfo())
		time.Sleep(time.Second * 3)
	}
}

func (l *LookHeapMemAndGoPrintInfo) run() {
	if l.running_go {
		return
	}
	if l.inlogger == nil {
		go LookHeapMemFunc(l)
	} else {
		go LookHeapMemFunc_log(l)
	}
	l.running_go = true
}

type LookHeapMemAndMaxMem struct {
	LookHeapMem
	running_go bool
	max        *runtime.MemStats
	lock_max   sync.Mutex
}

func NewLookHeapMemAndMaxMem() *LookHeapMemAndMaxMem {
	look := &LookHeapMemAndMaxMem{LookHeapMem: LookHeapMem{&runtime.MemStats{}}}
	look.run()
	return look
}

func LookHeapMem_maxmemFunc(look *LookHeapMemAndMaxMem) {
	for {
		mems := new(runtime.MemStats)
		runtime.ReadMemStats(mems)

		look.lock_max.Lock()
		if mems.HeapSys >= look.max.HeapSys && mems.HeapInuse > look.max.HeapInuse {
			look.max = mems
		}
		look.lock_max.Unlock()

		time.Sleep(time.Millisecond * 100)
	}
}

func (l *LookHeapMemAndMaxMem) run() {
	if l.running_go {
		return
	}
	if l.max == nil {
		mems := new(runtime.MemStats)
		runtime.ReadMemStats(mems)
		l.lock_max.Lock()
		l.max = mems
		l.lock_max.Unlock()
	}
	go LookHeapMem_maxmemFunc(l)

	l.running_go = true
}
func (l *LookHeapMemAndMaxMem) GetMaxMemInfo() runtime.MemStats {
	l.lock_max.Lock()
	defer l.lock_max.Unlock()

	return *l.max
}

func (l *LookHeapMemAndMaxMem) GetMaxHeapMemStatsMiB() string {
	mems := l.GetMaxMemInfo()
	return fmt.Sprintf("[Max] HeapSys: %dMiB, HeapIdle: %dMiB, HeapReleased: %dMiB, HeapInuse: %dMiB, HeapAlloc: %dMiB, HeapObjects: %d",
		mems.HeapSys/1024/1024, mems.HeapIdle/1024/1024, mems.HeapReleased/1024/1024, mems.HeapInuse/1024/1024, mems.HeapAlloc/1024/1024, mems.HeapObjects)
}
