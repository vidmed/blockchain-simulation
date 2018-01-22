package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vidmed/blockchain-simulation/simulator"
	"github.com/vidmed/logger"
)

var (
	configFileName = flag.String("config", "config.toml", "Config file name")
	sim            simulator.Simulator
)

func init() {
	flag.Parse()
	_, err := NewConfig(*configFileName)
	if err != nil {
		logger.Get().Fatalf("ERROR loading config: %s\n", err.Error())
	}
	// Init logging, logger goes first since other components may use it
	logger.Init(int(GetConfig().Main.LogLevel))
}

func main() {
	sim = simulator.NewSimulator(
		GetConfig().Main.FlushPeriod,
		GetConfig().Main.MaxTransactions,
		GetConfig().Main.FlushFile)
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
			sim.Close()
			logger.Get().Fatal(err.Error())
		}
	}()

	<-stop

	timeout := 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	logger.Get().Infof("Shutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Get().Errorf("Error: %v\n", err)
	} else {
		logger.Get().Infof("Server stopped")
	}
	cancel()

	// close simulator after server shutted down
	sim.Close()
}

func tx(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	k := q.Get("key")
	if k == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key is required"))
		return
	}
	v := q.Get("value")
	if v == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("value is required"))
		return
	}
	sim.Input() <- simulator.NewTransaction(k, v)
	w.WriteHeader(http.StatusOK)
}
