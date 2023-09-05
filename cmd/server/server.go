package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/NayronFerreira/client-server-api/internal/data"
	"github.com/NayronFerreira/client-server-api/internal/domain/entity"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CambioHandle)
	http.ListenAndServe(":8080", mux)
}

func CambioHandle(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/cotacao" {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	cambioResponse, err := GetCambioUSDToBRLWithContext(req.Context(), res)
	if err != nil {
		log.Println("Falha ao buscar CambioResponse:", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	database, err := data.CreateDBConnection()
	if err != nil {
		log.Print("Falha ao abrir coinexao com database", err)
	}
	defer database.Close()
	data.InsertCambioDB(database, cambioResponse)
	if err != nil {
		log.Print("Falha ao salvar cambio Response no database", err)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(&cambioResponse)
	res.WriteHeader(http.StatusOK)
}

func GetCambioUSDToBRLWithContext(ctxReq context.Context, res http.ResponseWriter) (*entity.CambioUSDBRL, error) {
	cambioResponse, err := getCambioUSDToBRL(ctxReq)
	if err != nil {
		return nil, err
	}
	return cambioResponse, nil
}

func getCambioUSDToBRL(ctxReq context.Context) (*entity.CambioUSDBRL, error) {
	ctx, cancel := context.WithTimeout(ctxReq, time.Millisecond*500)
	defer cancel()
	cambioReq, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	cambioRes, err := http.DefaultClient.Do(cambioReq)
	if err != nil {
		return nil, err
	}
	defer cambioRes.Body.Close()
	cambioResBody, err := io.ReadAll(cambioRes.Body)
	if err != nil {
		return nil, err
	}
	var cambioUSDBRL *entity.CambioUSDBRL
	err = json.Unmarshal(cambioResBody, &cambioUSDBRL)
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		log.Println("Solicitação processada com sucesso")
	case <-ctxReq.Done():
		log.Println("Solicitação cancelada pelo Client")
	}
	return cambioUSDBRL, nil
}
