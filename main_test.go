package main

import (
	"strconv"
	"testing"
	"time"
)

func TestParseAndStoreTime(t *testing.T) {
	testTime := time.Now()
	storedTime = &testTime
	storedTimeUnixString := strconv.FormatInt(storedTime.Unix(), 10)
	err := parseAndStoreTime(storedTimeUnixString)
	if err != nil {
		t.Error("error parsing string: " + err.Error())
	}

	if storedTime.Unix() != time.Now().Unix() {
		t.Error("wrong parse value", storedTime.Unix())
	}

}

func TestFailParseAndStoreTime(t *testing.T) {
	err := parseAndStoreTime("this is not a number")
	if err == nil {
		t.Error("we were expecting an error")
	}
}
