package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	ds := NewCacheDataStore()

	r := chi.NewRouter()
	r.Use(AuthMiddleware())
	r.Post("/note/{noteTitle}", newPostHandler(ds))
	r.Get("/note/{noteTitle}", newGetHandler(ds))
	r.Delete("/note/{noteTitle}", newDeleteHandler(ds))

	http.ListenAndServe(":8080", r)
}
