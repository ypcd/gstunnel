/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
*/
package timerm

import (
    "time"
)
type Timer struct {
    oldtime time.Time
    timeout time.Duration
}

func (t *Timer) Run() bool{
    if time.Since(t.oldtime).Seconds() > t.timeout.Seconds(){
        t.oldtime = time.Now()
        return true
    }
    return false
}
func (t *Timer) Boot() {
    t.oldtime = time.Now()
}
func CreateTimer(timeout time.Duration) Timer{
    t1 := Timer{time.Now(), timeout}
    return t1
}
