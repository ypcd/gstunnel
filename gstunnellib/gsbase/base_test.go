package gsbase

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func Test_race1(t *testing.T) {
	info, ok := debug.ReadBuildInfo()
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
