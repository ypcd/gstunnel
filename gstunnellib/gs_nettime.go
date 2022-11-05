package gstunnellib

import (
	"fmt"
	"time"
)

type NetTime interface {
	Add(time.Duration)
	PrintString() string
}

type netTimeImp struct {
	min, max, sum time.Duration
	count         uint64
	name          string
}

func NewNetTimeImp() NetTime {
	return &netTimeImp{}
}

func NewNetTimeImpName(name string) NetTime {
	return &netTimeImp{name: name}
}

func (nt *netTimeImp) Add(t time.Duration) {
	if nt.count == 0 {
		nt.min = t
		nt.max = t
		nt.sum += t
	} else {
		if t < nt.min {
			nt.min = t
		}
		if t > nt.max {
			nt.max = t
		}
		nt.sum += t
	}
	nt.count += 1
}

func (nt *netTimeImp) PrintString() string {
	if nt.count == 0 {
		return fmt.Sprintf("%s time: max:%s  min:%s  avg:%s  sum:%s  count:%d\n",
			nt.name, nt.max, nt.min, nt.sum, nt.sum, nt.count)
	}
	return fmt.Sprintf("%s time: max:%s  min:%s  avg:%s  sum:%s  count:%d\n",
		nt.name, nt.max, nt.min, nt.sum/time.Duration(nt.count), nt.sum, nt.count)
}
