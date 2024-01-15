package gstunnellib

import (
	"encoding/json"
	"fmt"
	"testing"
)

/*
func Test_gsstatus1(t *testing.T) {
	gid := NewGIdImp()
	s1 := NewGsStatusImp(gid)
	fmt.Println(string(s1.GetJson()))

	for i := 0; i < 1000*1000; i++ {
		id1 := gid.GenerateId()
		if s1.GetStatusData().Gid != id1 {
			panic("Error.")
		}
	}
}
*/

func Test_getnetconn(t *testing.T) {
	conn1, conn2 := GetRDNetConn()
	fmt.Println("conn:", conn1, conn2)
	conn3, conn4 := GetRDNetConn()
	fmt.Println("conn:", conn3, conn4)
}

func Test_statusConnListImp(t *testing.T) {
	cl := newStatusConnListImp()

	connList1, ok := cl.(*statusConnListImp)
	if !ok {
		t.Fatal("Error.")
	}
	conn1, conn2 := GetRDNetConn()
	connList1.Add(1, conn1, conn2)
	connList1.Add(2, conn2, conn1)

	j1, err := json.Marshal(connList1.ConnList)
	checkError_exit(err)
	fmt.Println(string(j1))

	connList1.Delete(1)
	fmt.Println(connList1.HTMLString())
}

func Test_GetStatusConnList(t *testing.T) {
	//g1 := NewGIdImp()
	v1 := NewGsStatusImp()
	cl := v1.GetStatusConnList()

	connList1, ok := cl.(*statusConnListImp)
	if !ok {
		t.Fatal("Error.")
	}
	conn1, conn2 := GetRDNetConn()
	connList1.Add(1, conn1, conn2)
	connList1.Add(2, conn2, conn1)

	j1, err := json.Marshal(connList1.ConnList)
	checkError_exit(err)
	fmt.Println(string(j1))

	connList1.Delete(1)
	fmt.Println(connList1.HTMLString())
}

func Test_GetStatusImp2(t *testing.T) {
	//g1 := NewGIdImp()
	v1 := NewGsStatusImp()
	//cl := v1.GetStatusConnList()
	j1, err := json.Marshal(v1)
	checkError_exit(err)
	fmt.Println(string(j1))

	cl := newStatusConnListImp()

	connList1, ok := cl.(*statusConnListImp)
	if !ok {
		t.Fatal("Error.")
	}
	conn1, conn2 := GetRDNetConn()
	connList1.Add(1, conn1, conn2)
	connList1.Add(2, conn2, conn1)

	j2, err := json.Marshal(connList1)
	checkError_exit(err)
	fmt.Println("statusConnListImp:", string(j2))

	j3, err := json.Marshal(cl)
	checkError_exit(err)
	fmt.Println("IStatusConnList:", string(j3))

}

func Test_GetStatusImp3(t *testing.T) {
	//g1 := NewGIdImp()
	v1 := NewGsStatusImp()
	//cl := v1.GetStatusConnList()

	conn1, conn2 := GetRDNetConn()
	v1.GetStatusConnList().Add(1, conn1, conn2)
	v1.GetStatusConnList().Add(2, conn2, conn1)

	j1, err := json.Marshal(v1)
	checkError_exit(err)
	fmt.Println(string(j1))

	connList1, ok := v1.GetStatusConnList().(*statusConnListImp)
	if !ok {
		t.Fatal("Error.")
	}

	j2, err := json.Marshal(connList1)
	checkError_exit(err)
	fmt.Println("statusConnListImp:", string(j2))

	j3, err := json.Marshal(v1.GetStatusConnList())
	checkError_exit(err)
	fmt.Println("IStatusConnList:", string(j3))

}

func Test_GetStatusConnListData(t *testing.T) {
	connList1 := newStatusConnListImp()

	conn1, conn2 := GetRDNetConn()
	connList1.Add(1, conn1, conn2)
	connList1.Add(2, conn2, conn1)

	scd := connList1.getStatusConnListData()
	_ = scd
}

func Test_gsstatus2(t *testing.T) {
	//gid := NewGIdImp()
	s1 := NewGsStatusImp()
	connList1 := s1.GetStatusConnList()
	conn1, conn2 := GetRDNetConn()
	fmt.Println(conn1, conn2)
	connList1.Add(1, conn1, conn2)
	connList1.Add(2, conn2, conn1)

	fmt.Println(string(s1.GetJson()))

}

func Test_gsstatus3(t *testing.T) {
	//gid := NewGIdImp()
	s1 := NewGsStatusImp()
	connList1 := s1.GetStatusConnList()
	conn1, conn2 := GetRDNetConn()
	fmt.Println(conn1, conn2)
	connList1.Add(1, conn1, conn2)
	connList1.Add(2, conn2, conn1)

	fmt.Println(string(s1.GetJson()))

	jstr, err := json.Marshal(s1)
	checkError_exit(err)
	fmt.Println(string(jstr))

}
