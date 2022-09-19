package gsbase

import (
	"bytes"
	"testing"
)

func Test_base1(t *testing.T) {
	buf := bytes.Buffer{}

	buf.Write([]byte("abc123"))
	buf.Reset()

	//buf.Read()
}
