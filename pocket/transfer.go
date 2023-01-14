package pocket

import (
	"database/sql"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (p *handler) Transfer(c echo.Context) error {
	return c.JSON(http.StatusOK, "hey")
}
