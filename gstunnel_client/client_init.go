package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gslog"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

const debug_gstClient bool = gsbase.Debug_gstClient

const g_version string = gstunnellib.G_Version

var g_Values = gstunnellib.NewGlobalValuesImp()

//var p = gstunnellib.Nullprint
//var pf = gstunnellib.Nullprintf

//var g_fpnull = os.DevNull

var g_key string
var g_key_rsa_server *gsrsa.RSA

var g_arg_key_gen = flag.String("g", "DefaultArg", fmt.Sprintf("'-g key'. Generate a random %d-byte base64 g_key.", gsbase.G_AesKeyLen))

var g_gsconfig *gstunnellib.GsConfig
var g_arg_gsconfig_path = flag.String("c", "config.client.json", "The gstunnel config file path.")

//var g_arg_gsversion = flag.String("v", g_version, "GST G_Version: "+g_version)

var g_arg_randLen = flag.Int("b", -1, "'-b number'. Generate a random number-byte []byte.")

var g_arg_gen_rsa_key = flag.String("rsa", "DefaultArg", "'-rsa key'. Generate a random 4096-bit rsa key and output it in base64 format.")

var g_goPackTotal, g_goUnpackTotal int32 = 0, 0

var g_logger *log.Logger

// Default 60s
var g_networkTimeout time.Duration

const g_net_read_size = 4 * 1024
const g_netPUn_chan_cache_size = 64

var g_Mt_model bool = true

// Default 6s
var g_tmr_display_time time.Duration

// Default 60s
var g_tmr_changekey_time time.Duration

//var bufPool sync.Pool

var g_RuntimeStatistics gstunnellib.IRuntime_statistics

var g_log_List gslog.Logger_List

var g_init_status = false

var g_gid = gstunnellib.NewGIdImp()

var g_gstst = gstunnellib.NewGsStatusImp()

var g_httpListenTryNum uint32 = 1000

func init_client_run() {

	if g_init_status {
		panic(errors.New("The init func is error."))
	} else {
		g_init_status = true
	}

	flag.Parse()
	if strings.EqualFold(*g_arg_key_gen, "key") {
		fmt.Println(gsrand.GetRDKeyBase64(gsbase.G_AesKeyLen))
		os.Exit(1)
	}
	/*
		 else if *g_arg_key_gen != "DefaultArg" {
			fmt.Println("Error: -g arg is error.\nCorrect arg: '-g g_key'.")
			os.Exit(-1)
		}
	*/
	if *g_arg_randLen > 0 {
		fmt.Println(gsrand.GetRDBytes(*g_arg_randLen))
		os.Exit(1)
	}
	if strings.EqualFold(*g_arg_gen_rsa_key, "key") {
		rsa1 := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)
		fmt.Printf("Private key:\n%s\n\nPublic key:\n%s\n\n", string(rsa1.PrivateKeyToBase64()), string(rsa1.PublicKeyToBase64()))
		os.Exit(1)
	}

	g_RuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	g_log_List.GenLogger = gslog.NewLoggerFileAndStdOut("gstunnel_client.log")
	g_log_List.GSIpLogger = gslog.NewLoggerFileAndStdOut("access.log")
	g_log_List.GSNetIOLen = gslog.NewLoggerFileAndLog("net_io_len.log", g_log_List.GenLogger.Writer())
	g_log_List.GSIpLogger.Println("Raw client access ip list:")

	g_logger = g_log_List.GenLogger

	g_logger.Println("gstunnel client.")
	g_logger.Println("VER:", g_version)

	g_gsconfig = gstunnellib.CreateGsconfig(*g_arg_gsconfig_path)
	g_Values.SetDebug(g_gsconfig.Debug)

	g_key = g_gsconfig.Key
	g_key_rsa_server = g_gsconfig.GetRSAServer()

	g_Mt_model = g_gsconfig.Mt_model

	g_tmr_display_time = time.Second * time.Duration(g_gsconfig.Tmr_display_time)
	g_tmr_changekey_time = time.Second * time.Duration(g_gsconfig.Tmr_changekey_time)
	g_networkTimeout = time.Second * time.Duration(g_gsconfig.NetworkTimeout)

	g_logger.Println("debug:", g_Values.GetDebug())

	g_logger.Printf("AES key len: %d Bit", gsbase.G_AESKeyBitLen)
	g_logger.Printf("RSA key len: %d Bit", gsbase.G_RSAKeyBitLen)

	g_logger.Println("g_Mt_model:", g_Mt_model)
	g_logger.Println("g_tmr_display_time:", g_tmr_display_time)
	g_logger.Println("g_tmr_changekey_time:", g_tmr_changekey_time)
	g_logger.Println("g_networkTimeout:", g_networkTimeout)

	g_logger.Println("info_protobuf:", gstunnellib.G_Info_protobuf)

	if g_Values.GetDebug() {
		var lock_httpListenAddr sync.Mutex
		httpListenAddr := ""
		var httpListenTryCount atomic.Uint32
		go func() {
			var err error
			for ; httpListenTryCount.Load() < g_httpListenTryNum; httpListenTryCount.Add(1) {
				lock_httpListenAddr.Lock()
				httpListenAddr = gstunnellib.GetNetLocalRDPort()
				lock_httpListenAddr.Unlock()
				err = http.ListenAndServe(httpListenAddr, nil)
			}
			panic(err)
		}()
		oldhttpListenTryCount := httpListenTryCount.Load()
		time.Sleep(500 * time.Millisecond)
		for httpListenTryCount.Load() > oldhttpListenTryCount {
			oldhttpListenTryCount = httpListenTryCount.Load()
			time.Sleep(100 * time.Millisecond)
		}
		lock_httpListenAddr.Lock()
		g_logger.Println("pprof server listen: " + httpListenAddr)
		lock_httpListenAddr.Unlock()
	}
	//g_Values.GetDebug() = false
	//go gstunnellib.RunGRuntimeStatistics_print(g_logger, g_RuntimeStatistics)
}

func init_client_test() {
	if g_init_status {
		panic(errors.New("The init func is error."))
	} else {
		g_init_status = true
	}

	g_RuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	g_log_List.GenLogger = gslog.NewLoggerFileAndStdOut("gstunnel_client.log")
	g_log_List.GSIpLogger = gslog.NewLoggerFileAndStdOut("access.log")
	g_log_List.GSNetIOLen = gslog.NewLoggerFileAndLog("net_io_len.log", g_log_List.GenLogger.Writer())
	g_log_List.GSIpLogger.Println("Raw client access ip list:")

	g_logger = g_log_List.GenLogger

	g_logger.Println("gstunnel client.")
	g_logger.Println("VER:", g_version)

	g_gsconfig = gstunnellib.CreateGsconfig(*g_arg_gsconfig_path)
	g_Values.SetDebug(g_gsconfig.Debug)

	g_key = g_gsconfig.Key
	g_key_rsa_server = g_gsconfig.GetRSAServer()

	g_Mt_model = g_gsconfig.Mt_model

	g_tmr_display_time = time.Second * time.Duration(g_gsconfig.Tmr_display_time)
	g_tmr_changekey_time = time.Second * time.Duration(g_gsconfig.Tmr_changekey_time)
	g_networkTimeout = time.Second * time.Duration(g_gsconfig.NetworkTimeout)

	g_logger.Println("debug:", g_Values.GetDebug())

	g_logger.Printf("AES key len: %d Bit", gsbase.G_AESKeyBitLen)
	g_logger.Printf("RSA key len: %d Bit", gsbase.G_RSAKeyBitLen)

	g_logger.Println("g_Mt_model:", g_Mt_model)
	g_logger.Println("g_tmr_display_time:", g_tmr_display_time)
	g_logger.Println("g_tmr_changekey_time:", g_tmr_changekey_time)
	g_logger.Println("g_networkTimeout:", g_networkTimeout)

	g_logger.Println("info_protobuf:", gstunnellib.G_Info_protobuf)

	if g_Values.GetDebug() {
		var lock_httpListenAddr sync.Mutex
		httpListenAddr := ""
		var httpListenTryCount atomic.Uint32
		go func() {
			var err error
			for ; httpListenTryCount.Load() < g_httpListenTryNum; httpListenTryCount.Add(1) {
				lock_httpListenAddr.Lock()
				httpListenAddr = gstunnellib.GetNetLocalRDPort()
				lock_httpListenAddr.Unlock()
				err = http.ListenAndServe(httpListenAddr, nil)
			}
			panic(err)
		}()
		oldhttpListenTryCount := httpListenTryCount.Load()
		time.Sleep(500 * time.Millisecond)
		for httpListenTryCount.Load() > oldhttpListenTryCount {
			oldhttpListenTryCount = httpListenTryCount.Load()
			time.Sleep(100 * time.Millisecond)
		}
		lock_httpListenAddr.Lock()
		g_logger.Println("pprof server listen: " + httpListenAddr)
		lock_httpListenAddr.Unlock()
	}
	//g_Values.GetDebug() = false
	//go gstunnellib.RunGRuntimeStatistics_print(g_logger, g_RuntimeStatistics)
}
