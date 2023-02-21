package utils

import (
	"reflect"
	"testing"
)

func TestISlice64toa(t *testing.T) {
	ids := []int64{3, 4, 5, 6, 1, 3, 4, 5, 56, 6, 213, 23, 1, 4, 32}
	toa := ISlice64toa(ids)
	if ok := reflect.DeepEqual(toa, "3,4,5,6,1,3,4,5,56,6,213,23,1,4,32"); !ok {
		t.Fail()
		t.Log("not compared!")
		return
	}
	t.Log(toa)
}
