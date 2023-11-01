package main

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/httpx"
	runtime2 "github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/example-domain/activity"
	"github.com/go-ai-agent/example-domain/google"
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
	displayRuntime()
	handler, status := startup(http.NewServeMux())
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

func displayRuntime() {
	fmt.Printf("addr : %v\n", addr)
	fmt.Printf("vers : %v\n", runtime.Version())
	fmt.Printf("os   : %v\n", runtime.GOOS)
	fmt.Printf("arch : %v\n", runtime.GOARCH)
	fmt.Printf("cpu  : %v\n", runtime.NumCPU())
}

func startup(r *http.ServeMux) (http.Handler, *runtime2.Status) {
	r.Handle(activity.HttpHandlerEndpoint, http.HandlerFunc(activity.HttpHandler))
	r.Handle(slo.HttpHandlerEndpoint, http.HandlerFunc(slo.HttpHandler))
	r.Handle(timeseries.HttpHandlerEndpoint, http.HandlerFunc(timeseries.HttpHandler))
	r.Handle(google.HttpHandlerEndpoint, http.HandlerFunc(google.HttpHandler))
	r.Handle(healthLivenessPattern, http.HandlerFunc(healthLivenessHandler))
	return r, runtime2.NewStatusOK()
}

func healthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime2.NewStatusOK()
	if status.OK() {
		httpx.WriteResponse[runtime2.LogError](w, []byte("up"), status, nil)
	} else {
		httpx.WriteResponse[runtime2.LogError](w, nil, status, nil)
	}
}
