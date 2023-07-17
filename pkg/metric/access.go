package metric

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"net/url"
)

func Get[E runtime.ErrorHandler](ctx context.Context, url *url.URL) ([]byte, *runtime.Status) {
	return nil, runtime.NewStatusOK()
}

func Post[E runtime.ErrorHandler](ctx context.Context, req *http.Request) *runtime.Status {
	return runtime.NewStatusOK()
}
