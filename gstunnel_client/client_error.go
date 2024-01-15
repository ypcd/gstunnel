package main

import (
	"errors"
	"strconv"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

func checkError(err error) {
	gserror.CheckErrorEx_exit(err, g_logger)
}

func checkError_NoExit(err error) {
	gserror.CheckErrorEx(err, g_logger)
}

// no exit
func checkError_info(err error) {
	gserror.CheckErrorEx_info(err, g_logger)
}

func checkError_panic(err error) {
	gserror.CheckErrorEx_panic(err)
}

func checkError_GSCtx(err error, gctx gstunnellib.IGSContext) {
	if err != nil {
		err2 := errors.New("[" + strconv.FormatUint(gctx.GetGsId(), 10) + "] " + err.Error())
		gserror.CheckErrorEx_exit(err2, g_logger)
	}
}

func checkError_NoExit_GSCtx(err error, gctx gstunnellib.IGSContext) {
	if err != nil {
		err2 := errors.New("[" + strconv.FormatUint(gctx.GetGsId(), 10) + "] " + err.Error())
		gserror.CheckErrorEx(err2, g_logger)
	}
}

func checkError_info_GSCtx(err error, gctx gstunnellib.IGSContext) {
	if err != nil {
		err2 := errors.New("[" + strconv.FormatUint(gctx.GetGsId(), 10) + "] " + err.Error())
		gserror.CheckErrorEx_info(err2, g_logger)
	}
}

func checkError_panic_GSCtx(err error, gctx gstunnellib.IGSContext) {
	if err != nil {
		err2 := errors.New("[" + strconv.FormatUint(gctx.GetGsId(), 10) + "] " + err.Error())
		gserror.CheckErrorEx_panic(err2)
	}
}
