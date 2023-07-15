package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

const g_version string = gstunnellib.G_Version

//gstunnellib.G_Version

var g_Values = gstunnellib.NewGlobalValuesImp()

//var p = gstunnellib.Nullprint
//var pf = gstunnellib.Nullprintf

var g_fpnull = os.DevNull

var g_key string

var g_arg_key_gen = flag.String("g", "DefaultArg", fmt.Sprintf("'-g key'. Generate a random %d-byte base64 key.", gsbase.G_AesKeyLen))

var g_gsconfig *gstunnellib.GsConfig
var g_arg_gsconfig_path = flag.String("c", "config.server.json", "The gstunnel config file path.")

//var g_arg_gsversion = flag.String("v", g_version, "GST G_Version: "+version)

var g_arg_randLen = flag.Int("b", -1, "'-b number'. Generate a random number-byte []byte.")

var g_Logger *log.Logger

const g_net_read_size = 4 * 1024
const g_netPUn_chan_cache_size = 64

var g_Mt_model bool = true

// Default 6s
var g_tmr_display_time time.Duration

// Default 60s
var g_tmr_changekey_time time.Duration

// Default 60s
var g_networkTimeout time.Duration

var g_RuntimeStatistics gstunnellib.Runtime_statistics

var g_log_List gstunnellib.Logger_List

var g_init_status = false

var g_gid = gstunnellib.NewGIdImp()

var g_gstst = gstunnellib.NewGsStatusImp(g_gid)

const g_connHandleGoNum = 16

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func init_server_run() {
	if g_init_status {
		panic(errors.New("the init func is error"))
	} else {
		g_init_status = true
	}
	/*
		for _, v := range os.Args {
			if v == "-g" {
			}
		}
	*/
	flag.Parse()
	if strings.EqualFold(*g_arg_key_gen, "key") {
		fmt.Println(gsrand.GetRDKeyBase64(gsbase.G_AesKeyLen))
		os.Exit(1)
	} else if *g_arg_key_gen != "DefaultArg" {
		fmt.Println("Error: -g arg is error.\nCorrect arg: '-g key'.")
		os.Exit(-1)
	}
	if *g_arg_randLen > 0 {
		fmt.Println(gsrand.GetRDBytes(*g_arg_randLen))
		os.Exit(1)
	}

	g_RuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	g_log_List.GenLogger = gstunnellib.NewLoggerFileAndStdOut("gstunnel_server.log")
	g_log_List.GSIpLogger = gstunnellib.NewLoggerFileAndStdOut("access.log")
	g_log_List.GSNetIOLen = gstunnellib.NewLoggerFileAndLog("net_io_len.log", g_log_List.GenLogger.Writer())
	g_log_List.GSIpLogger.Println("Gstunnel client access ip list:")

	g_Logger = g_log_List.GenLogger

	g_Logger.Println("gstunnel server.")

	g_Logger.Println("VER:", g_version)

	g_gsconfig = gstunnellib.CreateGsconfig(*g_arg_gsconfig_path)

	g_Values.SetDebug(g_gsconfig.Debug)

	g_Mt_model = g_gsconfig.Mt_model

	g_tmr_display_time = time.Second * time.Duration(g_gsconfig.Tmr_display_time)
	g_tmr_changekey_time = time.Second * time.Duration(g_gsconfig.Tmr_changekey_time)
	g_networkTimeout = time.Second * time.Duration(g_gsconfig.NetworkTimeout)

	g_key = g_gsconfig.Key

	g_Logger.Println("debug:", g_Values.GetDebug())

	g_Logger.Println("g_Mt_model:", g_Mt_model)
	g_Logger.Println("g_tmr_display_time:", g_tmr_display_time)
	g_Logger.Println("g_tmr_changekey_time:", g_tmr_changekey_time)
	g_Logger.Println("g_networkTimeout:", g_networkTimeout)

	g_Logger.Println("info_protobuf:", gstunnellib.G_Info_protobuf)

	if g_Values.GetDebug() {
		go func() {
			g_Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6070", nil))
		}()
		g_Logger.Println("pprof server listen: localhost:6070")
	}
	/*
		if g_gsconfig.WebUI {
			go gsweb.Run(g_gstst)
		}
	*/
	//g_Values.GetDebug() = false

	//go gstunnellib.RunGRuntimeStatistics_print(g_Logger, g_RuntimeStatistics)

}

func init_server_test() {
	if g_init_status {
		panic(errors.New("the init func is error"))
	} else {
		g_init_status = true
	}

	g_RuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	g_log_List.GenLogger = gstunnellib.NewLoggerFileAndStdOut("gstunnel_server.log")
	g_log_List.GSIpLogger = gstunnellib.NewLoggerFileAndStdOut("access.log")
	g_log_List.GSNetIOLen = gstunnellib.NewLoggerFileAndLog("net_io_len.log", g_log_List.GenLogger.Writer())
	g_log_List.GSIpLogger.Println("Gstunnel client access ip list:")

	g_Logger = g_log_List.GenLogger

	g_Logger.Println("gstunnel server.")

	g_Logger.Println("VER:", g_version)

	g_gsconfig = gstunnellib.CreateGsconfig(*g_arg_gsconfig_path)

	g_Values.SetDebug(g_gsconfig.Debug)

	g_Mt_model = g_gsconfig.Mt_model

	g_tmr_display_time = time.Second * time.Duration(g_gsconfig.Tmr_display_time)
	g_tmr_changekey_time = time.Second * time.Duration(g_gsconfig.Tmr_changekey_time)
	g_networkTimeout = time.Second * time.Duration(g_gsconfig.NetworkTimeout)

	g_key = g_gsconfig.Key

	g_Logger.Println("debug:", g_Values.GetDebug())

	g_Logger.Println("g_Mt_model:", g_Mt_model)
	g_Logger.Println("g_tmr_display_time:", g_tmr_display_time)
	g_Logger.Println("g_tmr_changekey_time:", g_tmr_changekey_time)
	g_Logger.Println("g_networkTimeout:", g_networkTimeout)

	g_Logger.Println("info_protobuf:", gstunnellib.G_Info_protobuf)

	if g_Values.GetDebug() {
		go func() {
			g_Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6070", nil))
		}()
		g_Logger.Println("pprof server listen: localhost:6070")
	}
	//g_Values.GetDebug() = false

	//go gstunnellib.RunGRuntimeStatistics_print(g_Logger, g_RuntimeStatistics)

}
