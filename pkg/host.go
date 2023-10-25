package main

/*
func Startup(mux *http.ServeMux) (http.Handler, *runtime.Status) {
	initMux(mux)
	return mux, runtime.NewStatusOK()
}

func initMux(r *http.ServeMux) {
	r.Handle(activity.EntryPath, http.HandlerFunc(activity.EntryHandler))
	r.Handle(slo.EntryPath, http.HandlerFunc(slo.EntryHandler))
	r.Handle(timeseries.EntryPath, http.HandlerFunc(timeseries.EntryHandler))
	r.Handle(healthLivenessPattern, http.HandlerFunc(HealthLivenessHandler))
}

func HealthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime.NewStatusOK()
	if status.OK() {
		httpx.WriteResponse[runtime.LogError](w, []byte("up"), status)
	} else {
		httpx.WriteResponse[runtime.LogError](w, nil, status)
	}
}


*/
