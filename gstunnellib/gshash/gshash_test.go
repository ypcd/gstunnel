package gshash

import "testing"

//AI generated
// Test_GetSha256Hex tests the GetSha256Hex function.
func Test_GetSha256Hex(t *testing.T) {
	data := []byte("Hello world")
	expected := "64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c"
	actual := GetSha256Hex(data)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
