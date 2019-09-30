package server

import (
	"encoding/json"
	"github.com/labbcb/brave/search"
	"github.com/labbcb/brave/variant"
	"log"
	"net/http"
)

func (s *Server) register() {
	s.Router.HandleFunc("/variants", s.adminOnly(s.handleInsertVariant())).Methods(http.MethodPost)
	s.Router.HandleFunc("/variants", s.adminOnly(s.handleRemoveVariants())).Methods(http.MethodDelete)
	s.Router.HandleFunc("/search", s.handleSearch()).Methods(http.MethodPost)
}

func (s *Server) handleInsertVariant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var v variant.Variant
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := s.InsertVariant(&v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]string{"id": v.ID}); err != nil {
			log.Println("encoding response to json:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input search.Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := s.Search(&input)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleRemoveVariants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		datasetID := r.FormValue("dataset")
		assemblyID := r.FormValue("assembly")

		if err := s.RemoveVariants(datasetID, assemblyID); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok || username != s.Username || password != s.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}
