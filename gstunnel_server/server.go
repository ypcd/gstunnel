/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

const version string = gstunnellib.Version

//gstunnellib.Version

var GValues = gstunnellib.NewGlobalValuesImp()

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string
var key_gen = flag.String("g", "1", "-g key. Generate a random 32-byte key.")

var gsconfig *gstunnellib.GsConfig
var gsconfig_path = flag.String("c", "config.server.json", "The gstunnel config file path.")

var Logger *log.Logger

const net_read_size = 4 * 1024
const netPUn_chan_cache_size = 64

var Mt_model bool = true

// defult 6s
var tmr_display_time time.Duration

// defult 60s
var tmr_changekey_time time.Duration

// defult 60s
var networkTimeout time.Duration

var GRuntimeStatistics gstunnellib.Runtime_statistics

var log_List gstunnellib.Logger_List

var init_status = false

var gid = gstunnellib.NewGIdImp()

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func init_server_run() {
	if init_status {
		panic(errors.New("The init func is error."))
	} else {
		init_status = true
	}
	/*
		for _, v := range os.Args {
			if v == "-g" {
			}
		}
	*/
	flag.Parse()
	if strings.Contains(strings.ToUpper(*key_gen), strings.ToUpper("key")) {
		fmt.Println(gstunnellib.GetRDKeyString32())
		os.Exit(1)
	}
	GRuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	log_List.GenLogger = gstunnellib.NewLoggerFileAndStdOut("gstunnel_server.log")
	log_List.GSIpLogger = gstunnellib.NewLoggerFileAndStdOut("access.log")
	log_List.GSNetIOLen = gstunnellib.NewLoggerFileAndLog("net_io_len.log", log_List.GenLogger.Writer())
	log_List.GSIpLogger.Println("Gstunnel client access ip list:")

	Logger = log_List.GenLogger

	Logger.Println("gstunnel server.")

	Logger.Println("VER:", version)

	gsconfig = gstunnellib.CreateGsconfig(*gsconfig_path)

	GValues.SetDebug(gsconfig.Debug)

	Mt_model = gsconfig.Mt_model

	tmr_display_time = time.Second * time.Duration(gsconfig.Tmr_display_time)
	tmr_changekey_time = time.Second * time.Duration(gsconfig.Tmr_changekey_time)
	networkTimeout = time.Second * time.Duration(gsconfig.NetworkTimeout)

	key = gsconfig.Key

	Logger.Println("debug:", GValues.GetDebug())

	Logger.Println("Mt_model:", Mt_model)
	Logger.Println("tmr_display_time:", tmr_display_time)
	Logger.Println("tmr_changekey_time:", tmr_changekey_time)
	Logger.Println("networkTimeout:", networkTimeout)

	Logger.Println("info_protobuf:", gstunnellib.Info_protobuf)

	if GValues.GetDebug() {
		go func() {
			Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6070", nil))
		}()
		Logger.Println("pprof server listen: localhost:6070")
	}
	//GValues.GetDebug() = false

	//go gstunnellib.RunGRuntimeStatistics_print(Logger, GRuntimeStatistics)

}

func init_server_test() {
	if init_status {
		panic(errors.New("The init func is error."))
	} else {
		init_status = true
	}

	//	flag.Parse()
	GRuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	log_List.GenLogger = gstunnellib.NewLoggerFileAndStdOut("gstunnel_server.log")
	log_List.GSIpLogger = gstunnellib.NewLoggerFileAndStdOut("access.log")
	log_List.GSNetIOLen = gstunnellib.NewLoggerFileAndLog("net_io_len.log", log_List.GenLogger.Writer())
	log_List.GSIpLogger.Println("Gstunnel client access ip list:")

	Logger = log_List.GenLogger

	Logger.Println("gstunnel server.")

	Logger.Println("VER:", version)

	gsconfig = gstunnellib.CreateGsconfig(*gsconfig_path)

	GValues.SetDebug(gsconfig.Debug)

	Mt_model = gsconfig.Mt_model

	tmr_display_time = time.Second * time.Duration(gsconfig.Tmr_display_time)
	tmr_changekey_time = time.Second * time.Duration(gsconfig.Tmr_changekey_time)
	networkTimeout = time.Second * time.Duration(gsconfig.NetworkTimeout)

	key = gsconfig.Key

	Logger.Println("debug:", GValues.GetDebug())

	Logger.Println("Mt_model:", Mt_model)
	Logger.Println("tmr_display_time:", tmr_display_time)
	Logger.Println("tmr_changekey_time:", tmr_changekey_time)
	Logger.Println("networkTimeout:", networkTimeout)

	Logger.Println("info_protobuf:", gstunnellib.Info_protobuf)

	if GValues.GetDebug() {
		go func() {
			Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6070", nil))
		}()
		Logger.Println("pprof server listen: localhost:6070")
	}
	//GValues.GetDebug() = false

	//go gstunnellib.RunGRuntimeStatistics_print(Logger, GRuntimeStatistics)

}

func main() {
	init_server_run()
	for {
		run()
	}
}

func run() {
	//defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	var lstnaddr, connaddr string

	lstnaddr = gsconfig.Listen
	connaddr = gsconfig.GetServer_rand()

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
			checkError_NoExit(err)
			continue
		}
		log_List.GSIpLogger.Printf("ip: %s\n", acc.RemoteAddr().String())

		service := connaddr

		dst, err := net.Dial("tcp", service)
		checkError(err)
		Logger.Println("conn.")

		gctx := gstunnellib.NewGsContextImp(gid.GetId())
		go srcTOdstUn(acc, dst, gctx)
		go srcTOdstP(dst, acc, gctx)
		Logger.Println("go.")
	}
}

func run_pipe_test_wg(dst net.Conn, acc net.Conn, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	Logger.Println("Test_Mt_model:", Mt_model)
	log_List.GSIpLogger.Printf("ip: %s\n", acc.RemoteAddr().String())

	gctx := gstunnellib.NewGsContextImp(gid.GetId())
	wg.Add(2)
	go srcTOdstUn_wg(acc, dst, wg, gctx)
	go srcTOdstP_wg(dst, acc, wg, gctx)
	Logger.Println("Gstunnel go.")

}

func IsTheVersionConsistent_send(dst net.Conn, apack gstunnellib.GsPack, wlent *int64) error {
	return gstunnellib.IsTheVersionConsistent_send(dst, apack, wlent)
}

func ChangeCryKey_send(dst net.Conn, apack gstunnellib.GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return gstunnellib.ChangeCryKey_send(dst, apack, ChangeCryKey_Total, wlent)
}

// service to gstunnel client
func srcTOdstP(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	if Mt_model {
		srcTOdstP_mt(src, dst, gctx)
	} else {
		srcTOdstP_st(src, dst, gctx)
	}
}

// gstunnel client to service
func srcTOdstUn(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	if Mt_model {
		srcTOdstUn_mt(src, dst, gctx)
	} else {
		srcTOdstUn_st(src, dst, gctx)
	}
}

// service to gstunnel client
func srcTOdstP_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg.Done()
	if Mt_model {
		srcTOdstP_mt(src, dst, gctx)
	} else {
		srcTOdstP_st(src, dst, gctx)
	}
}

// gstunnel client to service
func srcTOdstUn_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg.Done()
	if Mt_model {
		srcTOdstUn_mt(src, dst, gctx)
	} else {
		srcTOdstUn_st(src, dst, gctx)
	}
}
