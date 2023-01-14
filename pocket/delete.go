package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

const (
	cStmt = "delete from pockets where id = $1 RETURNING id;"
)

func (h handler) Delete(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	id := c.Param("id")
	var p pocket
	err := c.Bind(&p)

	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, id).Scan(&lastInsertId)
	if err != nil {
		logger.Error("delete row error", zap.Error(err))
		return err
	}

	logger.Info("delete successfully", zap.Int64("id", lastInsertId))
	return c.JSON(http.StatusAccepted, zap.Error(err))
}
