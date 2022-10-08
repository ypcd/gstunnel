/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package gstunnellib

import (
	"log"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
)

const Version string = gsbase.Version

var p func(...interface{}) (int, error)

var begin_Dinfo int = 0

var debug_tag bool

var commonIV = []byte{171, 158, 1, 73, 31, 98, 64, 85, 209, 217, 131, 150, 104, 219, 33, 220}

var logger *log.Logger

const Info_protobuf bool = true

const Deep_debug = gsbase.Deep_debug

var RunTimeDebugInfoV1 RunTimeDebugInfo

var key_defult string = "1234567890123456"

func Nullprint(v ...interface{}) (int, error)                       { return 1, nil }
func Nullprintf(format string, a ...interface{}) (n int, err error) { return 1, nil }

func init() {
	debug_tag = false
	p = Nullprint

	logger = NewFileLogger("gstunnellib.log")
	//debug_tag = true
	//p = fmt.Println

	RunTimeDebugInfoV1 = NewRunTimeDebugInfo()
}
