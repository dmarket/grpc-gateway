package prometheus

import "net/http"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func NewStatusRecorder(w http.ResponseWriter) *statusRecorder {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK
	return &statusRecorder{w, http.StatusOK}
}
