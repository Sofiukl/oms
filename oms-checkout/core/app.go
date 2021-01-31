package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/sofiukl/oms/oms-checkout/utils"
	"github.com/sofiukl/oms/oms-checkout/worker"
	"github.com/sofiukl/oms/oms-core/models"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
)

// WorkQueue - This is the work queues
var WorkQueue = make(chan worker.Work, 100)

// App - Application
type App struct {
	Router *mux.Router
	Conn   *pgxpool.Pool
	Config utils.Config
	Lock   *sync.RWMutex
}

// Initialize - This function initializes the application
func (a *App) Initialize() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.Connect(context.Background(), config.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	dbConnectMsg := fmt.Sprintf("Connected to DB %s", config.DBURL)
	fmt.Println(dbConnectMsg)
	a.Conn = conn
	a.Router = mux.NewRouter()
	a.Config = config
	a.Lock = &sync.RWMutex{}
	a.initializeRoutes()
}

// Run - This functio funs the application
func (a *App) Run(address string) {
	fmt.Println("Application is running on port", address)
	if err := http.ListenAndServe(address, a.Router); err != nil {
		log.Fatal(err)
	}
}

func (a *App) initializeRoutes() {
	s := a.Router.PathPrefix("/checkout-service/api/v1").Subrouter()
	s.HandleFunc("/checkout/", a.checkout).Methods("POST")
}

func (a *App) checkout(w http.ResponseWriter, r *http.Request) {
	body := parseBody(r)
	fmt.Println("IN THE STARTING OF HANDLER")
	work := worker.Work{Work: body, Config: a.Config, Conn: a.Conn, Lock: a.Lock}
	WorkQueue <- work
	fmt.Println("Work request queued")

	w.WriteHeader(http.StatusCreated)
	return

}

func parseBody(r *http.Request) models.CheckoutModel {
	decoder := json.NewDecoder(r.Body)

	var checkoutBody models.CheckoutModel
	err := decoder.Decode(&checkoutBody)

	if err != nil {
		log.Fatal(err)
	}
	return checkoutBody
}
