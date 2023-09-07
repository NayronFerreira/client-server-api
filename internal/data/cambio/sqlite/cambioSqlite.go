package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/NayronFerreira/client-server-api/internal/domain/entity"
)

const databaseFilePath string = "../../internal/data/cambio/sqlite/database.db"

func CreateDBConnection() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", databaseFilePath)
	if err != nil {
		return nil, err
	}
	return database, nil
}

func NewCambioTableIfNecessary(database *sql.DB) (*sql.DB, error) {
	_, err := database.Exec(queryCreateCambioTableIfNecessary())
	if err != nil {
		return nil, err
	}
	return database, nil
}

func InsertCambioDB(database *sql.DB, cambio *entity.CambioUSDBRL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	database, err := NewCambioTableIfNecessary(database)
	if err != nil {
		return err
	}
	stmt, err := database.PrepareContext(ctx, queryInsertCambio())
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,
		&cambio.USDBRL.Code, &cambio.USDBRL.Codein, &cambio.USDBRL.High, &cambio.USDBRL.Low,
		&cambio.USDBRL.VarBid, &cambio.USDBRL.PctChange, &cambio.USDBRL.Bid, &cambio.USDBRL.Ask,
		&cambio.USDBRL.Timestamp, &cambio.USDBRL.CreateDate)
	return err
}

func queryInsertCambio() string {
	return `
	INSERT INTO cambio_usdbrl (
		code, codein, high, low, var_bid, pct_change, bid, ask, timestamp, created_date
	) 
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)
`
}

func queryCreateCambioTableIfNecessary() string {
	return `
    CREATE TABLE IF NOT EXISTS cambio_usdbrl (
        id INTEGER PRIMARY KEY,
        code TEXT,
        codein TEXT,
        high TEXT,
        low TEXT,
        var_bid TEXT,
        pct_change TEXT,
        bid TEXT,
        ask TEXT,
        timestamp TEXT,
        created_date TEXT
    )
`
}
