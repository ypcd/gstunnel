package gstunnellib

import (
	"testing"
)

func Test_gscontext1(t *testing.T) {
	gc := NewGsContextImp(123)

	if gc.GetGsId() != 123 {
		panic("Error.")
	}
}
