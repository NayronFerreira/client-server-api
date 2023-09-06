package handle

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NayronFerreira/client-server-api/internal/data"
	"github.com/NayronFerreira/client-server-api/internal/service"
)

func CambioHandle(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/cotacao" {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	cambioResponse, err := service.GetCambioUSDToBRLWithReqContext(req.Context(), res)
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
