package main

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/future-architect/tagscanner/examples/restmap"
)

type Request struct {
	Method  string                `rest:"method"`
	Auth    string                `rest:"header:Authorization"`
	TraceID string                `rest:"cookie:trace-id"`
	Title   string                `rest:"body:title-field"`
	File    multipart.File        `rest:"body:file-field"`
	Header  *multipart.FileHeader `rest:"body:file-field"`
	Ctx     context.Context       `rest:"context"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/user/{userid}", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		restmap.Decode(&req, r)
		w.Write([]byte("welcome"))
	})
	fmt.Println("start server at :3000")
	http.ListenAndServe(":3000", r)
}
