package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func routes() *http.ServeMux {
	app := &PDFApp{
		MaxFileSize:  50 << 20,
		MaxTotalSize: 100 << 20,
		MaxMemory:    4 << 20,
	}
	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	mux.Handle("/internal/metrics", promhttp.Handler())
	mux.HandleFunc("/merge", securityMiddleware(app.MergeHandler))

	return mux
}
