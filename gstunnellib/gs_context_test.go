package gstunnellib

import (
	"testing"
)

func Test_gscontext1(t *testing.T) {
	gc := NewGSContextImp(123, NewGsStatusImp())

	if gc.GetGsId() != 123 {
		panic("Error.")
	}
}
