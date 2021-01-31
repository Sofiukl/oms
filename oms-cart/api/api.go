package api

import (
	"net/http"

	"github.com/sofiukl/oms/oms-cart/utils"
	"github.com/sofiukl/oms/oms-core/models"
)

// FindCart - find cart
func FindCart(id string, w http.ResponseWriter, r *http.Request) {

	cart := models.CartModel{
		ID: "c1",
		Products: []models.ProductModel{
			{ProdID: "p1", Qty: 1},
		},
	}

	// Return response
	utils.RespondWithJSON(w, http.StatusOK, "Product find successfully", "", cart)
}
