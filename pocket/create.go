package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	crestePkStmt = "INSERT INTO pockets(name, category, currency,balance) values($1, $2, $3, $4) RETURNING id;"
)

func (h *handler) CreatePocket(c echo.Context) error {

	logger := mlog.L(c)
	ctx := c.Request().Context()

	var pk pocket
	err := c.Bind(&pk)

	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, crestePkStmt, pk.Name, pk.Category, pk.Currency, pk.Balance).Scan(&lastInsertId)

	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))

	pk.ID = lastInsertId
	return c.JSON(http.StatusCreated, pk)

}
