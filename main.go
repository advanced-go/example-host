package main

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/handler"
	"github.com/advanced-go/core/http2"
	runtime2 "github.com/advanced-go/core/runtime"
	"github.com/advanced-go/example-domain/activity"
	"github.com/advanced-go/example-domain/google"
	"github.com/advanced-go/example-domain/slo"
	"github.com/advanced-go/example-domain/timeseries"
	"github.com/advanced-go/messaging/mux"
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
	setRuntimeEnvironment()
	setAccessLogging()

	displayRuntime()
	handler, status := startup(http.NewServeMux())
	if !status.OK() {
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("started : %v", time.Since(start)))

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
	fmt.Printf("addr    : %v\n", addr)
	fmt.Printf("vers    : %v\n", runtime.Version())
	fmt.Printf("os      : %v\n", runtime.GOOS)
	fmt.Printf("arch    : %v\n", runtime.GOARCH)
	fmt.Printf("cpu     : %v\n", runtime.NumCPU())
	fmt.Printf("env     : %v\n", runtime2.EnvStr())
}

func setRuntimeEnvironment() {
	// Set runtime environment
	//runtime2.SetTestEnvironment()
}

func setAccessLogging() {
	//SetAccessHandler(nil)
	access.EnableTestLogHandler()
	//access.EnableInternalLogging()
}

func startup(r *http.ServeMux) (http.Handler, runtime2.Status) {
	mux.Handle(activity.PkgPath, activity.HttpHandler)
	mux.Handle(slo.PkgPath, slo.HttpHandler)
	mux.Handle(timeseries.PkgPath, timeseries.HttpHandler)
	mux.Handle(google.PkgPath, google.HttpHandler)
	r.Handle(healthLivenessPattern, http.HandlerFunc(healthLivenessHandler))
	r.Handle("/", http.HandlerFunc(mux.HttpHandler))
	// Start agent
	//agent.Run(time.Second * 2)
	return handler.HttpHostMetricsHandler(r, ""), runtime2.NewStatusOK()
}

func healthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime2.NewStatusOK()
	if status.OK() {
		http2.WriteResponse[runtime2.LogError](w, []byte("up"), status, nil)
	} else {
		http2.WriteResponse[runtime2.LogError](w, nil, status, nil)
	}
}
