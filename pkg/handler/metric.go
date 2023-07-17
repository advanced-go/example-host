package handler

import (
	"github.com/go-ai-agent/core/exchange"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/example-host/pkg/metric"
	"net/http"
)

func MetricHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	status := runtime.NewHttpStatusCode(http.StatusBadRequest)

	if r.Method == http.MethodGet {
		buf, status = metric.Get[runtime.LogError](runtime.ContextWithRequest(r), r.URL)
	} else {
		status = metric.Post[runtime.LogError](runtime.ContextWithRequest(r), r)
	}
	exchange.WriteResponse[runtime.LogError](w, buf, status)
}
