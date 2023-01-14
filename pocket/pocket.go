package pocket

import (
	"database/sql"

	"github.com/kkgo-software-engineering/workshop/config"
)

type Pocket struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type Err struct {
	Message string
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

type Err struct {
	Message string `json:"message"`
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}
