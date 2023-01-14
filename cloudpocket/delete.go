package cloudpocket

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

const (
	cStmt = "delete tb where id = $1 RETURNING id;"
)

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (h handler) Delete(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	id := c.Param("id")
	// var e expense.Expenses
	// err := c.Bind(&e)
	// if err != nil {
	// 	logger.Error("bad request body", zap.Error(err))
	// 	return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	// }

	var lastInsertId int64
	err := h.db.QueryRowContext(ctx, cStmt, id).Scan(&lastInsertId)
	if err != nil {
		logger.Error("delete row error", zap.Error(err))
		return err
	}

	logger.Info("delete successfully", zap.Int64("id", lastInsertId))
	return c.JSON(http.StatusAccepted, zap.Error(err))
}
