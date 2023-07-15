package gstunnellib

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type GsStatus interface {
	GetJson() []byte
	GetStatusData() *gsStatusData
	GetStatusConnList() StatusConnList
}

type gsStatusImp struct {
	Gid             GId
	StatusConnList1 StatusConnList
}

type gsStatusData struct {
	Gid      uint64
	ConnList *statusConnListData
}

func NewGsStatusImp(ingid GId) GsStatus {
	return &gsStatusImp{Gid: ingid,
		StatusConnList1: newStatusConnListImp()}
}

func (s *gsStatusImp) GetJson() []byte {
	sd := s.GetStatusData()
	j1, err := json.Marshal(sd)
	CheckError_exit(err)
	return j1
}

func (s *gsStatusImp) GetStatusData() *gsStatusData {
	return &gsStatusData{Gid: s.Gid.GetId(), ConnList: s.GetStatusConnList().getStatusConnListData()}
}

func (s *gsStatusImp) GetStatusConnList() StatusConnList {
	return s.StatusConnList1
}

type StatusConnList interface {
	Add(id uint64, conn1, conn2 net.Conn)
	Delete(id uint64)
	HTMLString() string
	getStatusConnListData() *statusConnListData
}

type statusConnListData struct {
	Data map[uint64]string
}

type statusConnListImp struct {
	ConnList map[uint64]string
	lock     *sync.Mutex
}

func newStatusConnListImp() StatusConnList {
	return &statusConnListImp{ConnList: make(map[uint64]string), lock: &sync.Mutex{}}
}

func (sc *statusConnListImp) Add(id uint64, conn1, conn2 net.Conn) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	_, ok := sc.ConnList[id]
	if !ok {
		connstr := fmt.Sprintf("gstServer: %s  %s  gst--rawService: %s  %s",
			conn1.LocalAddr().String(), conn1.RemoteAddr().String(),
			conn2.LocalAddr().String(), conn2.RemoteAddr().String())
		sc.ConnList[id] = connstr
	}
}

func (sc *statusConnListImp) Delete(id uint64) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	delete(sc.ConnList, id)
}

func (sc *statusConnListImp) HTMLString() string {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	s1 := "len: %d<br>"
	outstr := fmt.Sprintf(s1, len(sc.ConnList))
	outList := make([]string, 1)

	for k, v := range sc.ConnList {
		outList = append(outList, fmt.Sprintf("id: %d : %s<br>", k, v))
	}

	sort.Slice(outList, func(i, j int) bool {
		defer func() {
			if msg := recover(); msg != nil {
				panic(msg)
			}
		}()

		if outList[i] == "" {
			return true
		}
		if outList[j] == "" {
			return false
		}

		v1 := outList[i]
		s := strings.Index(v1, ":")
		e := strings.Index(v1[s+2:], ":")
		if e == -1 {
			return true
		}
		vv1, err := strconv.Atoi(v1[s+2 : s+2+e-1])
		CheckError_exit(err)

		v1 = outList[j]
		s = strings.Index(v1, ":")
		e = strings.Index(v1[s+2:], ":")
		if e == -1 {
			return false
		}
		vv2, err := strconv.Atoi(v1[s+2 : s+2+e-1])
		CheckError_exit(err)
		return vv1 < vv2
	})
	outList[0] = outstr
	return strings.Join(outList, "")
}

func (sc *statusConnListImp) getStatusConnListData() *statusConnListData {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	d1 := make(map[uint64]string)
	for k := range sc.ConnList {
		d1[k] = sc.ConnList[k]
	}
	return &statusConnListData{d1}
}
