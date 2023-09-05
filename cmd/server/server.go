package main

import (
	"github.com/NayronFerreira/client-server-api/internal/controller"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	controller.NewCambioController()
}
