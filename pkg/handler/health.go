package handler

import (
	"github.com/go-ai-agent/core/exchange"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

func HealthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime.NewStatusOK()
	if status.OK() {
		exchange.WriteResponse[runtime.LogError](w, []byte("up"), status)
	} else {
		exchange.WriteResponse[runtime.LogError](w, nil, status)
	}
}
