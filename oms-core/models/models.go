package models

// GenericResponse - Common Response structure
type GenericResponse struct {
	Error   bool                   `json:"error"`
	Message string                 `json:"message,omitempty"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Details string                 `json:"details"`
}

// Product - model for Product
type Product struct {
	ID          string `json:"id,omitempty" bson:"id" mapstructure:"id"`
	Name        string `json:"name,omitempty" bson:"name" mapstructure:"name"`
	Description string `json:"description" bson:"description" mapstructure:"description"`
	AvailQty    int    `json:"avail_qty" bson:"avail_qty" mapstructure:"avail_qty" `
	ReserveQty  int    `json:"reserve_qty" bson:"reserve_qty" mapstructure:"reserve_qty"`
}

// ProductModel - model for taking product info
type ProductModel struct {
	ID  string `json:"id"`
	Qty int    `json:"quantity"`
}

// CartModel - cart model
type CartModel struct {
	ID       string         `json:"id"`
	Products []ProductModel `json:"products"`
}

// CheckoutModel - model for checkout
type CheckoutModel struct {
	//Product ProductModel `json:"product"`
	CartID string  `json:"cart_id"`
	Amount float64 `json:"amount"`
}
