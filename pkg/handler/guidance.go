package handler

import (
	"github.com/go-ai-agent/core/exchange"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/example-host/pkg/guidance"
	"net/http"
)

func GuidanceHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	status := runtime.NewHttpStatusCode(http.StatusBadRequest)

	if r.Method == http.MethodGet {
		buf, status = guidance.Get[runtime.LogError](runtime.ContextWithRequest(r), r.URL)
	} else {
		if r.Method == http.MethodPost {
			status = guidance.Post[runtime.LogError](runtime.ContextWithRequest(r), r)
		}
	}
	exchange.WriteResponse[runtime.LogError](w, buf, status)
}
