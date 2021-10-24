package gstunnellib

import "testing"

func Test_gorou_status(t *testing.T) {
	grs := CreateGorouStatus()
	grs.SetOk()
	ok := grs.IsOk()
	if !ok {
		t.Fatal("error.")
	}
	grs.SetClose()
	ok = grs.IsOk()
	if ok {
		t.Fatal("error.")
	}

}
