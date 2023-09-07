package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*600)
	defer cancel()

	bid, err := getBidCambio(ctx, "http://localhost:8080/cotacao")
	if err != nil {
		log.Fatal(err)
	}

	const cotacaoFilePath string = "../../internal/data/cambio/cotacao/cotaxao.txt"

	err = saveCambioToFile(bid, cotacaoFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Dólar: %s", bid)
}

func getBidCambio(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP Status: %s", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var cambio CambioUSDBRL
	err = json.Unmarshal(resBody, &cambio)
	if err != nil {
		return "", err
	}

	return cambio.USDBRL.Bid, nil
}

func saveCambioToFile(bid, filename string) error {
	arq, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer arq.Close()

	_, err = arq.WriteString("Dólar:" + bid)
	if err != nil {
		return err
	}

	return nil
}

type CambioUSDBRL struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}
