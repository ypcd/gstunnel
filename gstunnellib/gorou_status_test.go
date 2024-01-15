package gstunnellib

import (
	"net"
	"strings"
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
	conn1, conn2 := net.Pipe()
	grs := NewGorouStatusNetConn([]net.Conn{conn1, conn2})

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

	time.Sleep(time.Second * 1)

	_, err := conn1.Write([]byte("123456"))
	if !strings.Contains(err.Error(), "closed") {
		t.Fatal("Error.")
	}
	_, err = conn2.Write([]byte("123456"))
	if !strings.Contains(err.Error(), "closed") {
		t.Fatal("Error.")
	}
}

func Test_gorou_status_netconn(t *testing.T) {
	conn1, conn2 := net.Pipe()
	g1 := NewGorouStatusNetConn([]net.Conn{conn1, conn2})
	if !g1.IsOk() {
		t.Fatal("error.")
	}
	g1.SetClose()
	if g1.IsOk() {
		t.Fatal("error.")
	}

	_, err := conn1.Write([]byte("123456"))
	if !strings.Contains(err.Error(), "closed") {
		t.Fatal("Error.")
	}
	_, err = conn2.Write([]byte("123456"))
	if !strings.Contains(err.Error(), "closed") {
		t.Fatal("Error.")
	}

}
