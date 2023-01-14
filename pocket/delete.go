package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

const (
	dStmt = "delete from pockets where id = $1 RETURNING id;"
	sStmt = "select id, name, category, currency, balance from pockets where id = $1"
)

func (h handler) Delete(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	id := c.Param("id")
	var p Pocket
	err := c.Bind(&p)

	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, Err{Message: "bad request body"})
	}

	rows := h.db.QueryRow(sStmt, id)
	err = rows.Scan(&p.ID, &p.Name, &p.Category, &p.Currency, &p.Balance)
	if err != nil {
		return c.JSON(http.StatusNotFound, Err{Message: "Cloud pocket not found"})
	}

	if p.Balance > 0 {
		return c.JSON(http.StatusBadRequest,
			Err{Message: "Unable to delete this Cloud Pocket\n there is amount left in this Cloud Pocket, please move money out and try again"})
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, dStmt, id).Scan(&lastInsertId)
	if err != nil {
		logger.Error("delete row error", zap.Error(err))
		return err
	}

	logger.Info("delete successfully", zap.Int64("id", lastInsertId))
	return c.JSON(http.StatusAccepted, Err{Message: "Cloud pocket deleted successfully"})
}
