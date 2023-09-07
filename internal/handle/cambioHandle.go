package handle

import (
	"encoding/json"
	"log"
	"net/http"

	data "github.com/NayronFerreira/client-server-api/internal/data/cambio/sqlite"
	"github.com/NayronFerreira/client-server-api/internal/service"
)

const cotacaoPath = "/cotacao"

func CambioHandle(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != cotacaoPath {
		http.NotFound(res, req)
		return
	}

	cambioResponse, err := service.GetCambioUSDToBRLWithReqContext(req.Context(), res)
	if err != nil {
		log.Printf("Falha ao buscar CambioResponse: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	database, err := data.CreateDBConnection()
	if err != nil {
		log.Printf("Falha ao abrir conexao com database: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer database.Close()

	if err := data.InsertCambioDB(database, cambioResponse); err != nil {
		log.Printf("Falha ao salvar cambio Response no database: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(&cambioResponse); err != nil {
		log.Printf("Falha ao codificar resposta JSON: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
