package controller

import (
	"net/http"

	"github.com/NayronFerreira/client-server-api/internal/handle"
)

func NewServerCambioController() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handle.CambioHandle)
	http.ListenAndServe(":8080", mux)
}
