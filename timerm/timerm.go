/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package timerm

import (
	"fmt"
	"time"
)

type Timer struct {
	oldtime time.Time
	timeout time.Duration
}

func (t *Timer) Run() bool {
	if time.Since(t.oldtime).Seconds() > t.timeout.Seconds() {
		t.oldtime = time.Now()
		return true
	}
	return false
}
func (t *Timer) Boot() {
	t.oldtime = time.Now()
}

func CreateTimer(timeout time.Duration) Timer {
	t1 := Timer{time.Now(), timeout}
	return t1
}

type RecoTime struct {
	oldtime  time.Time
	min, max int64
	avg      int64
}

func (t *RecoTime) Run() int64 {
	t2 := time.Since(t.oldtime).Microseconds()

	if t.min == 0 {
		t.min = t2
	}
	if t.max == 0 {
		t.max = t2
	}
	if t2 < t.min {
		t.min = t2
	}
	if t2 > t.max {
		t.max = t2
	}

	t.avg = (t.avg + t2) / 2

	t.oldtime = time.Now()
	return t2
}

func (t *RecoTime) RunDisplay(sv1 string) {
	t2 := t.Run()
	fmt.Println(sv1, t2)
}

func (t *RecoTime) StringAll() string {
	return fmt.Sprintf("min: %d, max: %d, avg: %d", t.min, t.max, t.avg)
}

func CreateRecoTime() *RecoTime {
	return &RecoTime{oldtime: time.Now(), min: 0, max: 0, avg: 0}
}
