package handler

import (
	"github.com/go-ai-agent/core/exchange"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/example-domain/environment"
	"net/http"
)

func TraceHandler(w http.ResponseWriter, r *http.Request) {
	buf, status := environment.GetTrace[runtime.LogError](runtime.ContextWithRequest(r), r.URL)
	exchange.WriteResponse[runtime.LogError](w, buf, status)
}
