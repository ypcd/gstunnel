/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

const version string = gstunnellib.Version

//gstunnellib.Version

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string

var gsconfig *gstunnellib.GsConfig

var Logger *log.Logger

var debug_server bool = false

const net_read_size = 4 * 1024
const netPUn_chan_cache_size = 64

var Mt_model bool = true
var tmr_display_time = time.Second * 5
var tmr_changekey_time = time.Second * 60

var networkTimeout time.Duration = time.Minute * 1

var GRuntimeStatistics gstunnellib.Runtime_statistics

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func init() {

	GRuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	Logger = gstunnellib.CreateFileLogger("gstunnel_server.log")

	Logger.Println("gstunnel server.")

	Logger.Println("VER:", version)

	gsconfig = gstunnellib.CreateGsconfig("config.server.json")

	debug_server = gsconfig.Debug

	Mt_model = gsconfig.Mt_model

	tmr_display_time = time.Second * time.Duration(gsconfig.Tmr_display_time)
	tmr_changekey_time = time.Second * time.Duration(gsconfig.Tmr_changekey_time)

	key = gsconfig.Key

	Logger.Println("debug:", debug_server)

	Logger.Println("Mt_model:", Mt_model)
	Logger.Println("tmr_display_time:", tmr_display_time)
	Logger.Println("tmr_changekey_time:", tmr_changekey_time)

	Logger.Println("info_protobuf:", gstunnellib.Info_protobuf)

	if debug_server {
		go func() {
			Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6070", nil))
		}()
		Logger.Println("pprof server listen: localhost:6070")
	}
	debug_server = false

	go gstunnellib.RunGRuntimeStatistics_print(Logger, GRuntimeStatistics)

}

func main() {
	for {
		run()
	}
}

func run() {
	defer gstunnellib.Panic_Recover(Logger)

	var lstnaddr, connaddr string

	if len(os.Args) == 4 {
		lstnaddr = os.Args[1]
		connaddr = os.Args[2]
		key = os.Args[3]
	} else {
		lstnaddr = gsconfig.Listen
		connaddr = gsconfig.GetServer_rand()
		key = gsconfig.Key
	}

	fmt.Println("Listen_Addr:", lstnaddr)
	fmt.Println("Conn_Addr:", connaddr)
	fmt.Println("Begin......")

	service := lstnaddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		acc, err := listener.Accept()
		if err != nil {
			checkError_NoExit(err)
			continue
		}

		service := connaddr
		//tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		//fmt.Println(tcpAddr)
		//checkError(err)
		dst, err := net.Dial("tcp", service)
		checkError(err)
		fmt.Println("conn.")

		go srcTOdstUn(acc, dst)
		go srcTOdstP(dst, acc)
		fmt.Println("go.")
	}
}

func run_pipe_test(ss gstestpipe.Service_server, gsc gstestpipe.GsClient) {
	defer gstunnellib.Panic_Recover(Logger)

	acc := gsc.GetServerConn()
	dst := ss.GetClientConn()

	go srcTOdstUn(acc, dst)
	go srcTOdstP(dst, acc)
	Logger.Println("Gstunnel go.")

}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
}

func IsTheVersionConsistent_send(dst net.Conn, apack gstunnellib.GsPack, wlent *int64) error {
	return gstunnellib.IsTheVersionConsistent_send(dst, apack, wlent)
}

func ChangeCryKey_send(dst net.Conn, apack gstunnellib.GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return gstunnellib.ChangeCryKey_send(dst, apack, ChangeCryKey_Total, wlent)
}

// service to gstunnel client
func srcTOdstP(src net.Conn, dst net.Conn) {
	if Mt_model {
		srcTOdstP_mt(src, dst)
	} else {
		srcTOdstP_st(src, dst)
	}
}

// gstunnel client to service
func srcTOdstUn(src net.Conn, dst net.Conn) {
	if Mt_model {
		srcTOdstUn_mt(src, dst)
	} else {
		srcTOdstUn_st(src, dst)
	}
}

func checkError(err error) {
	gstunnellib.CheckErrorEx_exit(err, Logger)
}

func checkError_NoExit(err error) {
	gstunnellib.CheckErrorEx(err, Logger)
}
