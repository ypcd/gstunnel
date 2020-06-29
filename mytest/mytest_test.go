package mytest

import (
	//	"bytes"
	//	"encoding/json"
	//	"flag"
	//	"fmt"
	//	"gstunnellib"
	//	"log"
	//	"net"
	//"net/http"
	//_ "net/http/pprof"
	//	"os"
	//	"runtime"
	//	"runtime/pprof"
	//	"sync/atomic"
	"testing"
	"time"
	//	"timerm"
)

//var p = gstunnellib.Nullprint
//var pf = gstunnellib.Nullprintf

//var fpnull = os.DevNull

func sleep() {
	time.Sleep(time.Second * 2)
}

func run() {
	for i := 0; i < 1000000000; i++ {
		_ = i + 1
	}
}

func Test_test2(t *testing.T) {
	sleep()
}
