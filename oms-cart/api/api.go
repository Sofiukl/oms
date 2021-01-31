package api

import (
	"encoding/json"
	"net/http"

	"github.com/sofiukl/oms/oms-cart/utils"
	"github.com/sofiukl/oms/oms-core/models"
)

// FindCart - find cart
func FindCart(id string, w http.ResponseWriter, r *http.Request) {

	cart := models.CartModel{
		ID: "c1",
		Products: []models.ProductModel{
			{ID: "p1", Qty: 1},
		},
	}

	// Return response
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(cart)
	json.Unmarshal(inrec, &inInterface)

	utils.RespondWithJSON(w, http.StatusOK, "Cart find successfully", "", inInterface)
}
