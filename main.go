package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"
	"context"
	"syscall"
	"flag"
	"github.com/vidmed/logger"
	"fmt"
)

var (
	configFileName = flag.String("config", "config.toml", "Config file name")
)

func init() {
	flag.Parse()
	_, err := NewConfig(*configFileName)
	if err != nil {
		logger.Get().Fatalf("ERROR loading config: %s\n", err.Error())
	}
	// Init logging, logger goes first since other components may use it
	logger.Init(GetConfig().Main.LogLevel)
}

func main() {
	runServer()
}

func runServer() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/tx", tx)
	hs := &http.Server{Addr: GetConfig().Main.ListenStr, Handler: nil}

	go func() {
		logger.Get().Infof("Listening on http://%s\n", hs.Addr)

		if err := hs.ListenAndServe(); err != http.ErrServerClosed {
			logger.Get().Fatal(err.Error())
		}
	}()

	<-stop

	timeout := 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Get().Infof("Shutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Get().Errorf("Error: %v\n", err)
	} else {
		logger.Get().Infof("Server stopped")
	}
}

func tx(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	k, ok := q["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key is required"))
		return
	}
	v, ok := q["value"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("value is required"))
		return
	}
	w.Write([]byte(fmt.Sprintf("%v %v", k, v)))
}
