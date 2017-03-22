/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
*/
package timerm

import (
    "testing"
    "time"
)

func Test_CreateTimer(t *testing.T) {
    tmr := CreateTimer(time.Second*1)
    tmr.Run()
    t.Log("ok.")
}

func Test_Run(t *testing.T) {
    tmr := CreateTimer(time.Second*1)
    for !tmr.Run() {}
    t.Log("ok.")
}

func Test_Boot(t *testing.T) {
    tmr := CreateTimer(time.Second*2)
    tmr2 := CreateTimer(time.Second*1)
    oldt := time.Now()
    for !tmr.Run() {
        if tmr2.Run() {
            tmr.Boot()
            break
        }
    }
    for !tmr.Run() {}
    //t.Log(time.Since(oldt).Seconds())
    if int(time.Since(oldt).Seconds()) == 3 {
        t.Log("ok.")
    } else {
        t.Error()
    }
}


