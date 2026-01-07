package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFApp struct {
	MaxFileSize  int64
	MaxTotalSize int64
	MaxMemory    int64
}

func (app *PDFApp) MergeHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, app.MaxTotalSize+(1<<20))

	err := r.ParseMultipartForm(app.MaxMemory)
	if err != nil {
		http.Error(w, "Arquivos muito grandes", http.StatusRequestEntityTooLarge)
		return
	}

	files := r.MultipartForm.File["pdfs"]

	if len(files) == 0 {
		http.Error(w, "Nenhum arquivo enviado", http.StatusBadRequest)
		return
	}

	var totalSize int64
	for _, f := range files {
		if f.Size > app.MaxFileSize {
			http.Error(w, "Um dos arquivos excede o limite individual", http.StatusRequestEntityTooLarge)
			return
		}
		totalSize += f.Size
	}

	if totalSize > (app.MaxTotalSize) {
		http.Error(w, "O conjunto de arquivos excede o limite permitido", http.StatusRequestEntityTooLarge)
		return
	}

	readers := make([]io.ReadSeeker, len(files))
	var wg sync.WaitGroup

	for i, fHeader := range files {
		wg.Add(1)
		go func(index int, header *multipart.FileHeader) {
			defer wg.Done()

			f, _ := header.Open()
			defer f.Close()

			buf := new(bytes.Buffer)
			io.Copy(buf, f)
			readers[index] = bytes.NewReader(buf.Bytes())
		}(i, fHeader)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/pdf")

	conf := model.NewDefaultConfiguration()
	err = api.MergeRaw(readers, w, false, conf)

	if err != nil {
		http.Error(w, "erro ao concatenar arquivos", http.StatusInternalServerError)
		log.Println("erro ao concatenar arquivos:", err)
	}

	opsProcessed.Inc()
}

func securityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "http://localhost:8081" && origin != "https://pdf.andreramalho.tech" && origin != "https://www.pdf.andreramalho.tech" {
			http.Error(w, "Acesso proibido", http.StatusForbidden)
			log.Println("Acesso negado")
			return
		}
		next(w, r)
	}
}
