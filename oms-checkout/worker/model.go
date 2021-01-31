package worker

import (
	"sync"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/sofiukl/oms/oms-checkout/utils"
	"github.com/sofiukl/oms/oms-core/models"
)

// Work - This is WorkRequest model
type Work struct {
	Work   models.CheckoutModel
	Config utils.Config
	Conn   *pgxpool.Pool
	Lock   *sync.RWMutex
}
