package main

import (
	"encoding/json"
	"net/http"
)

func newServer(cfg Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg)
	})
func addDetails(){
mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg)
	})
	
}
	
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
}
