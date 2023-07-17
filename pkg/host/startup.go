package host

import (
	"github.com/go-ai-agent/core/resource"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

func Startup[E runtime.ErrorHandler, O runtime.OutputHandler](mux *http.ServeMux) (http.Handler, *runtime.Status) {
	var e E

	err := initLogging()
	if err != nil {
		return nil, e.Handle(nil, "/host/startup/logging", err)
	}
	initMux(mux)
	status := startupResources[E, O]()
	return mux, status
}

func Shutdown() {
	resource.Shutdown()
}

func startupResources[E runtime.ErrorHandler, O runtime.OutputHandler]() *runtime.Status {
	return resource.Startup[E, O](time.Second*5, nil)
}

func initLogging() error {
	return nil
}
