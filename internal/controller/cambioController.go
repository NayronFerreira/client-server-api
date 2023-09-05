package controller

import (
	"net/http"

	"github.com/NayronFerreira/client-server-api/internal/service"
)

func NewCambioController() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", service.CambioHandle)
	http.ListenAndServe(":8080", mux)
}
