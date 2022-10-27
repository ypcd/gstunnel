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
	"sync"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
	"github.com/ypcd/gstunnel/v6/timerm"
)

const version string = gstunnellib.Version

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string

var gsconfig *gstunnellib.GsConfig

var goPackTotal, goUnpackTotal int32 = 0, 0

var Logger *log.Logger

var networkTimeout time.Duration = time.Minute * 1

var debug_client bool = false

const net_read_size = 4 * 1024
const netPUn_chan_cache_size = 64

var Mt_model bool = true
var tmr_display_time = time.Second * 6
var tmr_changekey_time = time.Second * 60

//var bufPool sync.Pool

var GRuntimeStatistics gstunnellib.Runtime_statistics

var log_List gstunnellib.Logger_List

func init() {

	GRuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	log_List.GenLogger = gstunnellib.NewFileLogger("gstunnel_client.log")
	log_List.GSIpLogger = gstunnellib.NewFileLogger("access.log")
	log_List.GSNetIOLen = gstunnellib.NewFileLogger("net_io_len.log")
	log_List.GSIpLogger.Println("Raw client access ip list:")

	Logger = log_List.GenLogger

	Logger.Println("gstunnel client.")
	Logger.Println("VER:", version)

	gsconfig = gstunnellib.CreateGsconfig("config.client.json")
	debug_client = gsconfig.Debug

	key = gsconfig.Key

	Mt_model = gsconfig.Mt_model

	tmr_display_time = time.Second * time.Duration(gsconfig.Tmr_display_time)
	tmr_changekey_time = time.Second * time.Duration(gsconfig.Tmr_changekey_time)

	Logger.Println("debug:", debug_client)

	Logger.Println("Mt_model:", Mt_model)
	Logger.Println("tmr_display_time:", tmr_display_time)
	Logger.Println("tmr_changekey_time:", tmr_changekey_time)

	Logger.Println("info_protobuf:", gstunnellib.Info_protobuf)

	if debug_client {
		go func() {
			Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6060", nil))
		}()
		Logger.Println("Debug server listen: localhost:6060")
	}
	//debug_client = false
	//go gstunnellib.RunGRuntimeStatistics_print(Logger, GRuntimeStatistics)
}

func main() {
	for {
		run()
	}
}

func run() {
	//defer gstunnellib.Panic_Recover(Logger)

	var lstnaddr string
	var connaddr []string

	lstnaddr = gsconfig.Listen
	connaddr = gsconfig.GetServers()

	Logger.Println("Listen_Addr:", lstnaddr)
	Logger.Println("Conn_Addr:", connaddr)
	Logger.Println("Begin......")

	service := lstnaddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		acc, err := listener.Accept()
		if err != nil {
			Logger.Println("Error:", err)
			continue
		}
		log_List.GSIpLogger.Printf("ip: %s\n", acc.RemoteAddr().String())

		server_conn_error_total := 0
		tmr := timerm.CreateTimer(time.Second * 10)
		for {
			if tmr.Run() {
				Logger.Println("Error: server_conn_error_total: ", server_conn_error_total)
				server_conn_error_total = 0
				tmr.Boot()
			}
			service := gsconfig.GetServers()[0]

			dst, err := net.Dial("tcp", service)
			//checkError(err)
			if err != nil {
				if server_conn_error_total > 1000 {
					fmt.Fprintln(os.Stderr, "Error: server_conn_error_total > 10000")
					Logger.Fatalln("Error: server_conn_error_total > 10000")
				}
				server_conn_error_total++
				Logger.Println("Error: [net.Dial('tcp', service)]:", err.Error())
				fmt.Fprintln(os.Stderr, "Error: [net.Dial('tcp', service)]:", err.Error())

				continue
			}
			Logger.Println("conn.", service)

			//acc: 		src---client
			//dst: 		client---serever
			//pack: 	acc read recv, dst wirte send.
			//unpack:	dst read recv, acc wirte send.
			go srcTOdstP_count(acc, dst)
			go srcTOdstUn_count(dst, acc)
			break
		}
		Logger.Println("go.")
	}
}

func nouse_run_pipe_test(sc gstestpipe.RawdataPiPe, gss gstestpipe.GsPiPe) {
	//defer gstunnellib.Panic_Recover(Logger)

	acc := sc.GetServerConn()
	dst := gss.GetConn()

	Logger.Println("Test_Mt_model:", Mt_model)
	log_List.GSIpLogger.Printf("ip: %s\n", acc.RemoteAddr().String())

	go srcTOdstP_count(acc, dst)
	go srcTOdstUn_count(dst, acc)
	Logger.Println("Gstunnel go.")

}

func run_pipe_test_wg(sc gstestpipe.RawdataPiPe, gss gstestpipe.GsPiPe, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover(Logger)

	acc := sc.GetServerConn()
	dst := gss.GetConn()

	Logger.Println("Test_Mt_model:", Mt_model)
	log_List.GSIpLogger.Printf("ip: %s\n", acc.RemoteAddr().String())

	wg.Add(2)
	go srcTOdstP_wg(acc, dst, wg)
	go srcTOdstUn_wg(dst, acc, wg)
	Logger.Println("Gstunnel go.")

}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
}

func srcTOdstP_count(src net.Conn, dst net.Conn) {
	atomic.AddInt32(&goPackTotal, 1)
	srcTOdstP(src, dst)
	atomic.AddInt32(&goPackTotal, -1)

}

func srcTOdstUn_count(src net.Conn, dst net.Conn) {
	atomic.AddInt32(&goUnpackTotal, 1)
	srcTOdstUn(src, dst)
	atomic.AddInt32(&goUnpackTotal, -1)
}

func IsTheVersionConsistent_send(dst net.Conn, apack gstunnellib.GsPack, wlent *int64) error {
	return gstunnellib.IsTheVersionConsistent_send(dst, apack, wlent)
}

func ChangeCryKey_send(dst net.Conn, apack gstunnellib.GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return gstunnellib.ChangeCryKey_send(dst, apack, ChangeCryKey_Total, wlent)
}

func srcTOdstP(src net.Conn, dst net.Conn) {
	if Mt_model {
		srcTOdstP_mt(src, dst)
	} else {
		srcTOdstP_st(src, dst)
	}
}

func srcTOdstUn(src net.Conn, dst net.Conn) {
	if Mt_model {
		srcTOdstUn_mt(src, dst)
	} else {
		srcTOdstUn_st(src, dst)
	}
}

func srcTOdstP_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	if Mt_model {
		srcTOdstP_mt(src, dst)
	} else {
		srcTOdstP_st(src, dst)
	}
}

func srcTOdstUn_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
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

func checkError_info(err error) {
	gstunnellib.CheckErrorEx_info(err, Logger)
}

func checkError_panic(err error) {
	gstunnellib.CheckErrorEx_panic(err, Logger)
}
