package gsbase

import (
	"fmt"
	runtimeDebug "runtime/debug"
	"testing"
)

func Test_race1(t *testing.T) {
	info, ok := runtimeDebug.ReadBuildInfo()
	if !ok {
		t.Fatal("Error.")
	}
	for _, s := range info.Settings {
		if s.Key == "-race" {
			fmt.Printf("debug.ReadBuildInfo() --race=%s\n", s.Value)
			fmt.Println("race state:", GetRaceState())
			return
		}
	}
	fmt.Println("debug.ReadBuildInfo() -race=false")
	fmt.Println("race state:", GetRaceState())

}

func Test_print_debug_list(t *testing.T) {
	fmt.Println(Print_debug_list())
}
