package response

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseWriter) WriteHeader(code int) {
	r.StatusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	r.Body = b
	return r.ResponseWriter.Write(b)
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}
