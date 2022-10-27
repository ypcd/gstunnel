package gstunnellib

import (
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func Test_gorou_status(t *testing.T) {
	grs := NewGorouStatus()
	grs.SetOk()
	ok := grs.IsOk()
	if !ok {
		t.Fatal("error.")
	}
	grs.SetClose()
	ok = grs.IsOk()
	if ok {
		t.Fatal("error.")
	}

}

func Test_gorou_status_loop(t *testing.T) {
	grs := NewGorouStatus()

	loop_total := 10000 * 100
	go func() {
		for i := 0; i < loop_total; i++ {
			grs.IsOk()
		}
	}()
	go func() {
		for i := 0; i < loop_total; i++ {
			grs.SetOk()
		}
	}()
	go func() {
		for i := 0; i < loop_total; i++ {
			grs.SetClose()
		}
	}()
}

func Test_gorou_status_loop2(t *testing.T) {
	conn1, _ := net.Pipe()
	grs := NewGorouStatusNetConn(conn1)

	var go_ok int32 = 1
	defer atomic.SwapInt32(&go_ok, 0)

	go func() {
		for atomic.LoadInt32(&go_ok) == 1 {
			grs.IsOk()
		}
	}()
	go func() {
		for atomic.LoadInt32(&go_ok) == 1 {
			grs.SetOk()
		}
	}()
	go func() {
		for atomic.LoadInt32(&go_ok) == 1 {
			grs.SetClose()
		}
	}()

	time.Sleep(time.Second * 3)
}

func Test_gorou_status_netconn(t *testing.T) {
	conn1, _ := net.Pipe()
	g1 := NewGorouStatusNetConn(conn1)
	if !g1.IsOk() {
		t.Fatal("error.")
	}
	g1.SetClose()
	if g1.IsOk() {
		t.Fatal("error.")
	}
}
