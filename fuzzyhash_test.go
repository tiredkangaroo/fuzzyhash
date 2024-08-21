package fuzzyhash

import (
	"testing"
)

func handleErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestHash(t *testing.T) {
	got := Hash([]byte("tre concrete is unprocessable there!"), 1)
	got2 := Hash([]byte("the concrete is unprocessable therz!"), 1)
	result1 := Equal(got, got2)
	t.Log(result1)
	if !result1 {
		t.Error("result 1 is wrong.")
		t.Fail()
	}

	got3 := Hash([]byte("6865206b6e656c7420746f207468652067726f756e6420616e642070756c6c6564206f757420612072696e6720616e642073616964"), 1)
	got4 := Hash([]byte("6865206b6e656c7420746f2074686520666c6f6f7220616e642070756c6c6564206f757420612072696e6720616e642073616964"), 1)
	result2 := Equal(got3, got4)
	t.Log(result2)
	if result2 {
		t.Error("result 2 is wrong.")
		t.Fail()
	}

	got5 := Hash([]byte("6865206b6e656c7420746f207468652067726f756e6420616e642070756c6c6564206f757420612072696e6720616e642073616964"), 2) // increased tolerance
	got6 := Hash([]byte("6865206b6e656c7420746f2074686520666c6f6f7220616e642070756c6c6564206f757420612072696e6720616e642073616964"), 2)   // increased tolerance
	result3 := Equal(got5, got6)
	t.Log(result3)
	if !result3 {
		t.Error("result 3 is wrong.")
		t.Fail()
	}

	got7 := Hash([]byte("elephants like horses"), 1)
	got8 := Hash([]byte("elephants are horses"), 1)
	result4 := Equal(got7, got8)
	t.Log(result4)
	if result4 {
		t.Error("result 4 is wrong.")
		t.Fail()
	}

	got9 := Hash([]byte("elephants like horses"), 1)
	got10 := Hash([]byte("elephants are like horses"), 1)
	result5 := Equal(got9, got10)
	t.Log(result5)
	if !result5 {
		t.Error("result 4 is wrong.")
		t.Fail()
	}
}
