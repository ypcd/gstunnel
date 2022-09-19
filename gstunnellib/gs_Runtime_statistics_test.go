package gstunnellib

import (
	"fmt"
	"testing"
	"time"
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
	//v1.AddTotalNetworkData(1)
	go func() {
		for {
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
