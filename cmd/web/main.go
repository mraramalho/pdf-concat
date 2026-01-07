package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mraramalho/toolkit"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pdf_merge_total_operations",
		Help: "O número total de merges realizados com sucesso",
	})
)

func main() {
	srv := &http.Server{
		Addr:    ":8081",
		Handler: routes(),
	}

	t := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 50,
		AllowedFileTypes: []string{"application/pdf"},
	}

	err := t.RunServer(context.Background(), srv, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

}
