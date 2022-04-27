package main

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

func newPostHandler(ds DataStore) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		note, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error reading note"))
			return 
		}
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("error retrieving user"))
			return 
		}
		ds.Store(user, chi.URLParam(r, "noteTitle"), string(note))
		w.Write([]byte("note stored successfully"))
	}
}

func newGetHandler(ds DataStore) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("error retrieving user"))
			return 
		}
		note, err := ds.Get(user, chi.URLParam(r, "noteTitle"))

		if err != nil {
			if _, ok := err.(NoteNotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("note not found"))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}
			return
		}  
		
		w.Write([]byte(note))
	}
}

func newDeleteHandler(ds DataStore) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("error retrieving user"))
			return 
		}
		err := ds.Delete(user, chi.URLParam(r, "noteTitle"))

		if err != nil {
			if _, ok := err.(NoteNotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("note not found"))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}
			return
		}  
		
		w.Write([]byte("note deleted"))
	}
}



var numLetterUnderscoreRegex = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func AuthMiddleware() func (http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// this is implementation of authorization is obviously terrible but the idea of some process checking header and "validating" user identity for downstream processing is what I'm after
			user := r.Header.Get("Authorization")

			if !numLetterUnderscoreRegex.Match([]byte(user)) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Invalid credentials - Please send a username containing numbers, letters, or underscores in the Authorization header"))
				return
			}

			h.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
		})
	}
}