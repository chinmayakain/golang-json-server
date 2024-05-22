package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
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
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (as *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/account", makeHTTPHandleFunc(as.handleGetAllAccount))
	mux.HandleFunc("GET /api/account/{id}", makeHTTPHandleFunc(as.handleGetAccountById))
	mux.HandleFunc("POST /api/account", makeHTTPHandleFunc(as.handleCreateAccount))
	mux.HandleFunc("DELETE /api/account/{id}", makeHTTPHandleFunc(as.handleDeleteAccount))
	mux.HandleFunc("PUT /api/account/transfer", makeHTTPHandleFunc(as.handleTransfer))

	log.Println("JSON API server listening on port: ", as.listenAddr)
	http.ListenAndServe(as.listenAddr, mux)
}

func (as *APIServer) handleGetAllAccount(w http.ResponseWriter, r *http.Request) error {
	var accounts []*Account
	var err error

	if accounts, err = as.store.GetAllAccounts(); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, nil)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (as *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}

	account, err := as.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (as *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccRequest); err != nil {
		return err
	}

	defer r.Body.Close()
	account := NewAccount(createAccRequest.FirstName, createAccRequest.LastName)

	if err := as.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (as *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}

	if err := as.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted account with account id": id})
}

func (as *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferReq)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func getId(r *http.Request) (int, error) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id: %v", idStr)
	}

	return id, nil
}
