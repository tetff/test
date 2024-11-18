package main

import (
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	ch         chan reqRes
	wg         sync.WaitGroup
	storedTime *time.Time
)

type reqRes struct {
	req *http.Request
	res http.ResponseWriter
}

func main() {

	ch = make(chan reqRes)
	storedTime = &time.Time{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", requestHandler)

	println("Starting service. Listening on localhost:51104")
	http.ListenAndServe(":51104", mux)
}

func requestHandler(res http.ResponseWriter, req *http.Request) {
	go processRequest()
	if req.Header.Get("Content-Type") == "text/plain" {
		println("sending request to processing")
		res.Header().Add("Content-Type", "text/plain")
		ch <- reqRes{req, res}
	} else {
		res.WriteHeader(http.StatusBadRequest)
		res.Header().Add("Content-Type", "text/plain")
		res.Write([]byte("content-type is not supported"))
	}
}

func processRequest() {
	reqRes := <-ch
	switch reqRes.req.Method {
	case http.MethodPost:
		println("recieving timestamp")
		reqBody, err := readBody(reqRes.res, reqRes.req)
		if err != nil {
			println("error reading timestamp: " + err.Error())
			reqRes.res.WriteHeader(http.StatusBadRequest)
			reqRes.res.Write([]byte("error reading request body: " + err.Error()))
			return
		}
		println("parsing timestamp")
		err = parseAndStoreTime(reqBody)
		if err != nil {
			println("error paring and storing timestamp: " + err.Error())
			reqRes.res.WriteHeader(http.StatusBadRequest)
			reqRes.res.Write([]byte("error parsing request body: " + err.Error()))
			return
		}
		println("sending OK response")
		reqRes.res.WriteHeader(http.StatusOK)
		reqRes.res.Write([]byte("unix timestamp successfuly stored"))
		return

	case http.MethodGet:
		println("returning stored timestamp")
		reqRes.res.WriteHeader(http.StatusOK)
		reqRes.res.Write([]byte(strconv.FormatInt(storedTime.Unix(), 10)))
		return

	default:
		println("recieved uncompatible request method")
		reqRes.res.WriteHeader(http.StatusBadRequest)
		reqRes.res.Write([]byte("request method is not supported"))
		return
	}

}

func readBody(res http.ResponseWriter, req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseAndStoreTime(t string) error {
	timeInt, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return err
	}
	myTime := time.Unix(timeInt, 0)
	storedTime = &myTime
	return nil
}
