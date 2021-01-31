package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pgx "github.com/jackc/pgx/v4"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"

	"github.com/sofiukl/oms/oms-core/models"
	"github.com/sofiukl/oms/oms-product/utils"
)

const (
	findProductQry = "select name, avail_qty, reserve_qty from product where id=$1 FOR UPDATE"
)

// FindProduct - find productc by id
func FindProduct(conn1 *pgxpool.Pool, config utils.Config, prodID string, w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), config.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	var name string
	var availQty int
	var reserveQty int

	// begin transaction eparate this as util
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), findProductQry, prodID).Scan(&name, &availQty, &reserveQty)
	if err != nil {
		log.Printf("FindProduct QueryRow failed: %v\n", err)
		utils.RespondWithError(w, 500, "Internal server error. Please try after some time.", "")
	}
	log.Println(name, availQty, reserveQty)
	prod := &models.Product{ID: prodID, Name: name, AvailQty: availQty, ReserveQty: reserveQty}

	// commit transaction separate this as common func
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Return response
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(prod)
	json.Unmarshal(inrec, &inInterface)

	utils.RespondWithJSON(w, http.StatusOK, "Product find successfully", "", inInterface)
}
