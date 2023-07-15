/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package timerm

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_CreateTimer(t *testing.T) {
	tmr := CreateTimer(time.Second * 1)
	tmr.Run()
	t.Log("ok.")
}

func Test_Run(t *testing.T) {
	tmr := CreateTimer(time.Second * 1)
	for !tmr.Run() {
	}
	t.Log("ok.")
}

func Test_Boot(t *testing.T) {
	tmr := CreateTimer(time.Second * 2)
	tmr2 := CreateTimer(time.Second * 1)
	oldt := time.Now()
	for !tmr.Run() {
		if tmr2.Run() {
			tmr.Boot()
			break
		}
	}
	for !tmr.Run() {
	}
	//t.Log(time.Since(oldt).Seconds())
	if int(time.Since(oldt).Seconds()) == 3 {
		t.Log("ok.")
	} else {
		t.Error()
	}
}

// Test_CreateRecoTime tests the CreateRecoTime function.
func Test_RecoTime(t *testing.T) {
	tmr := CreateRecoTime()
	for i := 0; i < 40; i++ {
		ts1 := time.Millisecond * time.Duration(15+rand.Float64()*10)
		t1 := time.Now()
		time.Sleep(ts1)
		fmt.Println("set sleep time:", ts1)
		fmt.Println("run sleep time:", time.Since(t1))
		tmr.RunDisplay("test")
		//fmt.Println(tmr.StringAll())
	}
	t.Log(tmr)
	fmt.Println(tmr.StringAll())

}

func Test_RecoTime2(t *testing.T) {
	tmr := CreateRecoTime()
	for i := 0; i < 10; i++ {
		tmr.RunDisplay("test")
		time.Sleep(time.Millisecond * 10)
		fmt.Println(tmr.StringAll())
		tmr.RunDisplay("test2")
		fmt.Println(tmr.StringAll())
	}
	t.Log(tmr)
	fmt.Println(tmr.StringAll())
}

func Test_RecoTime3(t *testing.T) {
	tmr := CreateRecoTime()

	tls := []time.Duration{600, 60, 6}

	for _, v := range tls {
		t1 := time.Now()
		time.Sleep(v * time.Millisecond)
		fmt.Println("set sleep time:", v*time.Millisecond)
		fmt.Println("run sleep time:", time.Since(t1))
		tmr.RunDisplay("test")
		//fmt.Println(tmr.StringAll())
	}
	t.Log(tmr)
	fmt.Println(tmr.StringAll())

}
