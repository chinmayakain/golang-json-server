package main

import (
	"encoding/json"
	"log"
	"net/http"
)


type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (as *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/accounts", makeHTTPHandleFunc(as.handleGetAllAccount))
	mux.HandleFunc("GET /api/accounts/{id}", makeHTTPHandleFunc(as.handleGetAccount))
	mux.HandleFunc("POST /api/accounts", makeHTTPHandleFunc(as.handleCreateAccount))
	mux.HandleFunc("DELETE /api/accounts", makeHTTPHandleFunc(as.handleDeleteAccount))
	mux.HandleFunc("PUT /api/accounts", makeHTTPHandleFunc(as.handleTransfer))

	log.Println("JSON API server listening on port: ", as.listenAddr)
	http.ListenAndServe(as.listenAddr, mux)
}

func (as *APIServer) handleGetAllAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Jan", "Gould")
	return WriteJSON(w, http.StatusOK, account)
}

func (as *APIServer)  handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	return WriteJSON(w, http.StatusOK, id)
}

func (as *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (as *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (as *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}