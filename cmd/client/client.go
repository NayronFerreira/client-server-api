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
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal(res.Status)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resBody))
	defer res.Body.Close()

	var bid CambioUSDBRL
	err = json.Unmarshal(resBody, &bid)
	if err != nil {
		panic(err)
	}
	insertCambioDolarInTheFile(bid.USDBRL.Bid)

}

func insertCambioDolarInTheFile(bid string) {
	arq, err := os.Create("cotacao.txt")
	if err != nil {
		panic(arq)
	}
	defer arq.Close()
	arq.WriteString("Dólar:" + bid)
}

type CambioUSDBRL struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}
