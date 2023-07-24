package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	go func() {
		for {
			if catchSignals(os.Interrupt, syscall.SIGTERM) {
				log.Println("shutting down...")
				<-time.After(2 * time.Second)
				if err := srv.Close(); err != nil {
					log.Printf("server close error: %s\n", err)
				}

				return
			}
		}
	}()

	log.Printf("starting api server. pid: %s\n", pid())

	if err := srv.ListenAndServe(); err != nil {
		switch {
		case errors.Is(http.ErrServerClosed, err):
			log.Println("gracefully shutdown")

		default:
			panic(err)
		}
	}
}

func catchSignals(ss ...os.Signal) bool {
	ntfy := make(chan os.Signal, 1)
	signal.Notify(ntfy, ss...)

	sig, ok := <-ntfy
	if !ok {
		log.Printf("signal channel unreadable")

		return false
	}

	log.Printf("got signal: %s\n", sig)

	return true
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
