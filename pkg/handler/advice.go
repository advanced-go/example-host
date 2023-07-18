package handler

import (
	"github.com/go-ai-agent/core/exchange"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/example-domain/actions"
	"net/http"
)

func AdviceHandler(w http.ResponseWriter, r *http.Request) {
	buf, status := actions.GetAdvice[runtime.LogError](runtime.ContextWithRequest(r), r.URL)
	exchange.WriteResponse[runtime.LogError](w, buf, status)
}
