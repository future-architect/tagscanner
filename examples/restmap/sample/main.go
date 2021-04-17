package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/osaki-lab/tagscanner/examples/restmap"
	"gitlab.com/osaki-lab/tagscanner/examples/restmap/sample"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/user/{userid}", func(w http.ResponseWriter, r *http.Request) {
		var req sample.GetUserRequest
		restmap.Decode(&req, r)
		w.Write([]byte("welcome"))
	})
	fmt.Println("start server at :3000")
	http.ListenAndServe(":3000", r)
}

