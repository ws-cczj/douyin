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

func TestSearchZero(t *testing.T) {
	ary := []int64{1, 3, 5, 6, 213, 45, 3, 5, 234, 1, 3, 0, 0, 0, 0, 0, 0}
	ary2 := []int64{1, 3, 5, 6, 213, 45, 3, 5, 234, 1, 3}
	ary3 := []int64{0}
	ary4 := []int64{}
	idx := SearchZero(ary)
	if idx != 11 {
		t.Log("not compared")
		t.FailNow()
	}
	t.Logf("split ary from zero left, res: %v", ary[:idx])
	idx2 := SearchZero(ary2)
	if idx2 != 11 {
		t.Log("not compared")
		t.FailNow()
	}
	t.Logf("split ary from zero left, res: %v", ary[:idx2])
	idx3 := SearchZero(ary3)
	if idx3 != 0 {
		t.Log("not compared")
		t.FailNow()
	}
	t.Logf("split ary from zero left, res: %v", ary[:idx3])
	idx4 := SearchZero(ary4)
	if idx4 != 0 {
		t.Log("not compared")
		t.FailNow()
	}
	t.Logf("split ary from zero left, res: %v", ary[:idx4])
}
