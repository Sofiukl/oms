package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	pgx "github.com/jackc/pgx/v4"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"

	"github.com/sofiukl/oms/oms-checkout/utils"
	"github.com/sofiukl/oms/oms-core/models"

	"github.com/mitchellh/mapstructure"
)

const (
	findProductQry = "select name, avail_qty, reserve_qty from product where id=$1 FOR UPDATE"
)

// CheckoutProduct - checkout the product
func CheckoutProduct(conn *pgxpool.Pool, config utils.Config, body models.CheckoutModel, lock *sync.RWMutex) {

	product, _ := findCartDetails(body.CartID)
	amount := body.Amount

	// begin transaction eparate this as util
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback(context.Background())

	msg, checkoutErr := checkout(tx, product, amount, lock)
	fmt.Println(msg)
	if checkoutErr != nil {
		// commit transaction separate this as common func
		fmt.Println(checkoutErr)
		err = tx.Commit(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		//utils.RespondWithError(w, http.StatusBadRequest, msg, "")
	}

	// commit transaction separate this as common func
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// Return response
	//utils.RespondWithJSON(w, http.StatusOK, msg, "", map[string]interface{}{"referenceNo": "R001"})
}

func updateReserveQty(tx pgx.Tx, prodID string, qty int, opType string) error {
	var qry string
	if opType == "increase" {
		qry = fmt.Sprintf("update product set reserve_qty = reserve_qty + %d where id =$1", qty)
	} else {
		qry = fmt.Sprintf("update product set reserve_qty = reserve_qty - %d where id =$1", qty)
	}

	//fmt.Println(qry)
	_, err := tx.Exec(context.Background(), qry, prodID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return fmt.Errorf("update quantity failed")
	}

	return nil
}

func updateAvailQty(tx pgx.Tx, prodID string, qty int) error {
	qry := fmt.Sprintf("update product set avail_qty = avail_qty - %d where id =$1", qty)
	_, err := tx.Exec(context.Background(), qry, prodID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return fmt.Errorf("update avail quantity failed")
	}

	return nil
}

func checkout(tx pgx.Tx, product models.ProductModel, amount float64, lock *sync.RWMutex) (string, error) {

	id := product.ID
	qty := product.Qty

	fmt.Println(id, qty)
	lock.Lock()
	defer lock.Unlock()
	prod, err := findProduct(id)

	fmt.Printf("%+v", prod)

	// var name string
	// var availQty int
	// var reserveQty int

	// err := tx.QueryRow(context.Background(), findProductQry, id).Scan(&name, &availQty, &reserveQty)
	// if err != nil {
	// 	log.Printf("FindProduct QueryRow failed: %v\n", err)
	// }
	// log.Println(name, availQty, reserveQty)
	// prod := &models.Product{ID: id, Name: name, AvailQty: availQty, ReserveQty: reserveQty}

	if err == nil {
		if prod.AvailQty-prod.ReserveQty < qty {
			fmt.Println("out of stock")
			//fmt.Printf("%+v", prod)
			return "The product is out of stock at this moment", nil
		} else {
			// process for checkout
			updateReserveQty(tx, id, qty, "increase")
			err = payment(id, amount)
			if err == nil {
				fmt.Println(err)
				fmt.Println(prod)
				updateAvailQty(tx, id, qty)
				updateReserveQty(tx, id, qty, "decrease")
				return "Yup! you successfully bought the product", nil
			} else {
				fmt.Println("2")
				prod, err = findProduct(id)
				updateReserveQty(tx, id, qty, "decrease")
				//fmt.Printf("%+v", prod)
				return "Your payment is not successfull", fmt.Errorf("Your payment is not successfull")
			}
		}
	} else {
		log.Println(err)
		return "Fail to checkout at this moment", fmt.Errorf("Fail to checkout st this moment")
	}
}

func payment(prodID string, amount float64) error {
	time.Sleep(2000)
	return nil
}

func findProduct(id string) (*models.Product, error) {
	var p models.Product
	var gresp models.GenericResponse
	link := fmt.Sprintf("http://localhost:3004/product/api/v1/find/%s", id)
	//fmt.Println(link)
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal("Fail to fetch product details from product service")
		return &p, fmt.Errorf("Fail to fetch produt")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	if err = json.Unmarshal(body, &gresp); err != nil {
		log.Println(err)
		return &p, fmt.Errorf("Fail to fetch produt")
	}
	//fmt.Println(gresp.Result)
	mapstructure.Decode(gresp.Result, &p)
	//fmt.Println(p)
	return &p, nil
}

func findCartDetails(id string) (models.ProductModel, error) {
	var gresp models.GenericResponse
	var cart models.CartModel
	var prod models.ProductModel

	link := fmt.Sprintf("http://localhost:3006/cart/api/v1/find/%s", id)
	//fmt.Println(link)
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal("Fail to fetch product details from product service")
		return prod, fmt.Errorf("Fail to fetch produt")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err = json.Unmarshal(body, &gresp); err != nil {
		log.Println(err)
		return prod, fmt.Errorf("Fail to fetch produt")
	}
	mapstructure.Decode(gresp.Result, &cart)
	fmt.Println(cart.Products)
	return cart.Products[0], nil
}
