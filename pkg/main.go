package main

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/httpx"
	runtime2 "github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/example-domain/activity"
	"github.com/go-ai-agent/example-domain/slo"
	"github.com/go-ai-agent/example-domain/timeseries"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const (
	addr                  = "0.0.0.0:8080"
	writeTimeout          = time.Second * 300
	readTimeout           = time.Second * 15
	idleTimeout           = time.Second * 60
	healthLivenessPattern = "/health/liveness"
)

func main() {
	start := time.Now()
	displayRuntime[runtime2.StdOutput]()
	handler, status := Startup(http.NewServeMux())
	if !status.OK() {
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("host started: %v", time.Since(start)))

	srv := http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		} else {
			log.Printf("HTTP server Shutdown")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func displayRuntime[O runtime2.OutputHandler]() {
	var o O
	o.Write(fmt.Sprintf("addr : %v", addr))
	o.Write(fmt.Sprintf("vers : %v", runtime.Version()))
	o.Write(fmt.Sprintf("os   : %v", runtime.GOOS))
	o.Write(fmt.Sprintf("arch : %v", runtime.GOARCH))
	o.Write(fmt.Sprintf("cpu  : %v", runtime.NumCPU()))
}

func Startup(mux *http.ServeMux) (http.Handler, *runtime2.Status) {
	initMux(mux)
	return mux, runtime2.NewStatusOK()
}

func initMux(r *http.ServeMux) {
	path := activity.EntryPath
	fmt.Printf("path: %v", path)
	r.Handle(activity.EntryPath, http.HandlerFunc(activity.EntryHandler))
	r.Handle(slo.EntryPath, http.HandlerFunc(slo.EntryHandler))
	r.Handle(timeseries.EntryPath, http.HandlerFunc(timeseries.EntryHandler))
	r.Handle(healthLivenessPattern, http.HandlerFunc(HealthLivenessHandler))
}

func HealthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime2.NewStatusOK()
	if status.OK() {
		httpx.WriteResponse[runtime2.LogError](w, []byte("up"), status)
	} else {
		httpx.WriteMinResponse[runtime2.LogError](w, status)
	}
}
