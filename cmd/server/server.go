package main

import (
	"net/http"

	"github.com/NayronFerreira/client-server-api/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", service.CambioHandle)
	http.ListenAndServe(":8080", mux)
}
