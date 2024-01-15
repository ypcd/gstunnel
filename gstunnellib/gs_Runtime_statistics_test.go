package gstunnellib

import (
	"fmt"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_Runtime_statistics_imp(t *testing.T) {
	v1 := NewRuntimeStatistics()
	//v1.AddTotalNetworkData(1)
	v1.AddServerTotalNetData_recv(1)
	v1.AddServerTotalNetData_send(1)
	v1.AddSrcTotalNetData_recv(1)
	v1.AddSrcTotalNetData_send(1)

	re, err := v1.GetJson()
	CheckError_test(err, t)
	fmt.Println(re)
	fmt.Println(string(re))
}

func Test_Runtime_statistics_imp_mgo(t *testing.T) {
	v1 := NewRuntimeStatistics()

	gos := NewGorouStatus()
	defer gos.SetClose()
	go func() {
		for gos.IsOk() {
			time.Sleep(time.Millisecond * 100)
			re, err := v1.GetJson()
			CheckError_test(err, t)
			_ = re
			fmt.Println(string(re))
		}
	}()

	for i := 0; i < 10000*100; i++ {
		v1.AddServerTotalNetData_recv(1)
		v1.AddServerTotalNetData_send(1)
		v1.AddSrcTotalNetData_recv(1)
		v1.AddSrcTotalNetData_send(1)
	}
}

func Test_Runtime_statistics_mgo(t *testing.T) {
	v1 := NewRuntimeStatistics()
	gos := NewGorouStatus()
	defer gos.SetClose()

	go func() {
		for gos.IsOk() {
			v1.AddServerTotalNetData_recv(int(gsrand.GetRDInt64()))
		}
	}()

	go func() {
		for gos.IsOk() {
			v1.AddServerTotalNetData_send(int(gsrand.GetRDInt64()))
		}
	}()

	go func() {
		for gos.IsOk() {
			v1.AddSrcTotalNetData_recv(int(gsrand.GetRDInt64()))
		}
	}()

	go func() {
		for gos.IsOk() {
			v1.AddSrcTotalNetData_send(int(gsrand.GetRDInt64()))
		}
	}()

	go func() {
		for gos.IsOk() {
			v1.GetJson()
		}
	}()

	time.Sleep(time.Second * 3)
}
