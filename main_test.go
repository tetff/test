package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParseAndStoreTime(t *testing.T) {
	success(t)
}

func TestParseAndStoreTimeStress(t *testing.T) {
	for i := 0; i < 10000; i++ {
		go func() {
			success(t)
		}()
	}
}

func TestFailParseAndStoreTime(t *testing.T) {
	fail(t)
}

func TestAlternatingSuccessAndFail(t *testing.T) {
	success(t)
	fail(t)
	success(t)
	success(t)
	fail(t)
	success(t)
	success(t)
	fail(t)
	success(t)
}

func success(t *testing.T) {
	reqBody := strings.NewReader(strconv.FormatInt(time.Now().Unix(), 10))
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:51104", reqBody)
	req.Header.Add("Content-Type", "text/plain")
	if err != nil {
		t.Error("error building request:", err)
	}

	res := httptest.NewRecorder()
	requestHandler(res, req)

	if res.Code != http.StatusOK {
		t.Error("wrong response code:", res.Code)
	}
}

func fail(t *testing.T) {
	reqBody := strings.NewReader("This is not a number")
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:51104", reqBody)
	req.Header.Add("Content-Type", "text/plain")
	if err != nil {
		t.Error("error building request:", err)
	}

	res := httptest.NewRecorder()
	requestHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Error("wrong response code:", res.Code)
	}
}
