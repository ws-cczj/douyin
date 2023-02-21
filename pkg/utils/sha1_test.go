package utils

import (
	"reflect"
	"testing"
)

// 40bd001563085fc35165329ea1ff5c5ecbdbbeef
func TestSHA1(t *testing.T) {
	sha1 := SHA1("123")
	if ok := reflect.DeepEqual(sha1, "40bd001563085fc35165329ea1ff5c5ecbdbbeef"); !ok {
		t.Fail()
		t.Logf("SHA1 exec res : %v", sha1)
	}
	t.Logf("SHA1 res: %v", sha1)
	t.Log("SHA1 test success!")
}
