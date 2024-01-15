package gsbase

import "fmt"

func Print_debug_list() string {
	return fmt.Sprintf(
		"gsbase debug_list:  Debug--%v  Debug_gstlib--%v  Debug_gstServer--%v  Debug_gstClient--%v\n",
		Debug, Debug_gstlib, Debug_gstServer, Debug_gstClient)
}
