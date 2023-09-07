package service

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/NayronFerreira/client-server-api/internal/constants"
	"github.com/NayronFerreira/client-server-api/internal/domain/entity"
)

func GetCambioUSDToBRLWithReqContext(ctxReq context.Context, res http.ResponseWriter) (*entity.CambioUSDBRL, error) {
	cambioResponse, err := getCambioUSDToBRL(ctxReq)
	if err != nil {
		return nil, err
	}
	return cambioResponse, nil
}

func getCambioUSDToBRL(ctxReq context.Context) (*entity.CambioUSDBRL, error) {
	ctx, cancel := context.WithTimeout(ctxReq, time.Millisecond*550)
	defer cancel()
	cambioReq, err := http.NewRequestWithContext(ctx, "GET", constants.URL_GET_CAMBIO, nil)
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
