package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

var PDFHandlerTestCases = []struct {
	testName           string
	filePaths          []string
	expectedStatusCode int
	expectsError       bool
	maxFileSize        int64
	maxTotalSize       int64
	maxMemory          int64
}{
	{
		testName:           "Sucesso com dois arquivos pequenos",
		filePaths:          []string{"./test_data/over5mb.pdf", "./test_data/smallFile.pdf"},
		expectedStatusCode: http.StatusOK,
		expectsError:       false,
		maxFileSize:        10 << 20,
		maxTotalSize:       20 << 20,
		maxMemory:          4 << 20,
	},
	{
		testName:           "Erro quando o um arquivo excede o limite",
		filePaths:          []string{"./test_data/over5mb.pdf", "./test_data/smallFile.pdf"},
		expectedStatusCode: http.StatusRequestEntityTooLarge,
		expectsError:       true,
		maxFileSize:        1 << 20,
		maxTotalSize:       10 << 20,
		maxMemory:          4 << 20,
	},
	{
		testName:           "Erro quando o total excede o limite total",
		filePaths:          []string{"./test_data/over5mb.pdf", "./test_data/over5mb.pdf"},
		expectedStatusCode: http.StatusRequestEntityTooLarge,
		expectsError:       true,
		maxFileSize:        6 << 20,
		maxTotalSize:       9 << 20,
		maxMemory:          4 << 20,
	},
	{
		testName:           "sucesso com arquivos grandes",
		filePaths:          []string{"./test_data/over5mb.pdf", "./test_data/over5mb.pdf"},
		expectedStatusCode: http.StatusOK,
		expectsError:       false,
		maxFileSize:        6 << 20,
		maxTotalSize:       20 << 20,
		maxMemory:          4 << 20,
	},
}

func TestHandlers_MergeHandler(t *testing.T) {
	for _, e := range PDFHandlerTestCases {
		t.Run(e.testName, func(t *testing.T) {
			app := &PDFApp{
				MaxFileSize:  e.maxFileSize,
				MaxTotalSize: e.maxTotalSize,
				MaxMemory:    e.maxMemory,
			}

			pr, pw := io.Pipe()
			writer := multipart.NewWriter(pw)
			wg := sync.WaitGroup{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				// A ordem LIFO garante que o writer feche (enviando boundary) antes do pipe
				defer pw.Close()
				defer writer.Close()

				for _, p := range e.filePaths {
					part, err := writer.CreateFormFile("pdfs", "teste.pdf")
					if err != nil {
						if !e.expectsError {
							t.Errorf("creating part form file error: %v", err)
						}
						return
					}

					f, err := os.Open(p)
					if err != nil {
						if !e.expectsError {
							t.Errorf("opening file error: %v", err)
						}
						return
					}

					_, err = io.Copy(part, f)
					f.Close()

					if err != nil {
						// Se o pipe fechou, significa que o handler parou de ler (esperado em erros de limite)
						if err != io.ErrClosedPipe && !e.expectsError {
							t.Errorf("copying file error: %v", err)
						}
						return
					}
				}
			}()

			req, err := http.NewRequest("POST", "/teste", pr)
			if err != nil {
				t.Fatalf("creating request error: %v", err)
			}
			req.Header.Add("Content-Type", writer.FormDataContentType())

			rr := httptest.NewRecorder()

			app.MergeHandler(rr, req)

			pr.Close()
			wg.Wait()

			if rr.Result().StatusCode != e.expectedStatusCode {
				t.Errorf("expected code %d, found %d", e.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
