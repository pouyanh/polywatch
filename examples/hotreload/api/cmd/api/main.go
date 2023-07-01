package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealthCheck)
	mux.HandleFunc("/", handleNow)
	srv := http.Server{
		Addr:    "",
		Handler: mux,
	}

	println("starting api server")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

type RspHealthCheck struct {
	Pid    string `json:"pid"`
	Status bool   `json:"status"`
}

func handleHealthCheck(rw http.ResponseWriter, _ *http.Request) {
	writeJson(rw, RspHealthCheck{Pid: pid(), Status: true})
}

type RspNow struct {
	Pid  string `json:"pid"`
	When string `json:"when"`
}

func handleNow(rw http.ResponseWriter, _ *http.Request) {
	writeJson(rw, RspNow{Pid: pid(), When: time.Now().Format(time.RFC1123)})
}

func writeJson(rw http.ResponseWriter, rsp any) {
	bb, err := json.Marshal(rsp)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(fmt.Sprint(err)))

		return
	}

	_, _ = rw.Write(bb)
}

func pid() string {
	return strconv.Itoa(os.Getpid())
}
