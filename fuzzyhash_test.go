package fuzzyhash

import (
	"reflect"
	"testing"
)

func handleErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestHash256(t *testing.T) {
	got := MustHash(Hash256([]byte("tre concrete is unprocessable there!"), 1))
	got2 := MustHash(Hash256([]byte("the concrete is unprocessable therz!"), 1))
	result1 := Equal(got, got2)
	t.Log(result1)
	if !result1 {
		t.Error("result 1 is wrong.")
		t.Fail()
	}

	got3 := MustHash(Hash256([]byte("6865206b6e656c7420746f207468652067726f756e6420616e642070756c6c6564206f757420612072696e6720616e642073616964"), 16))
	got4 := MustHash(Hash256([]byte("6865206b6e656c7420746f2074686520666c6f6f7220616e642070756c6c6564206f757420612072696e6720616e642073616964"), 16))
	result2 := !Equal(got3, got4)
	t.Log(result2)
	if !result2 {
		t.Error("result 2 is wrong.")
		t.Fail()
	}

	got5 := MustHash(Hash256([]byte("6865206b6e656c7420746f207468652067726f756e6420616e642070756c6c6564206f757420612072696e6720616e642073616964"), 8)) // increased tolerance
	got6 := MustHash(Hash256([]byte("6865206b6e656c7420746f2074686520666c6f6f7220616e642070756c6c6564206f757420612072696e6720616e642073616964"), 8))   // increased tolerance
	result3 := Equal(got5, got6)
	t.Log(result3)
	if !result3 {
		t.Error("result 3 is wrong.")
		t.Fail()
	}

	got7 := MustHash(Hash256([]byte("elephants like horses"), 16))
	got8 := MustHash(Hash256([]byte("elephants are horses"), 16))
	result4 := !Equal(got7, got8)
	t.Log(result4)
	if !result4 {
		t.Error("result 4 is wrong.")
		t.Fail()
	}

	got9 := MustHash(Hash256([]byte("elephants like horses"), 1))
	got10 := MustHash(Hash256([]byte("elephants are like horses"), 1))
	result5 := Equal(got9, got10)
	t.Log(result5)
	if !result5 {
		t.Error("result 5 is wrong.")
		t.Fail()
	}
}

func TestMarshal(t *testing.T) {
	got11 := MustHash(Hash256([]byte("elephants are like horses"), 16)).MarshalString()
	result6 := got11 == "045827d5d10a1300430e95683ed7b57b295aae163b789a24f246f33a446a57c6:16"
	t.Log(result6)
	if !result6 {
		t.Error("result 6 is wrong.")
		t.Fail()
	}
}

func TestUnmarshal(t *testing.T) {
	got12, err := UnmarshalString("045827d5d10a1300430e95683ed7b57b295aae163b789a24f246f33a446a57c6:16")
	if err != nil {
		t.Errorf("unexpected error getting result 7: %s", err.Error())
	}
	result7 := reflect.DeepEqual(*got12, *MustHash(Hash256([]byte("elephants are like horses"), 16)))
	t.Log(result7)
	if !result7 {
		t.Error("result 7 is wrong.")
		t.Fail()
	}
}
