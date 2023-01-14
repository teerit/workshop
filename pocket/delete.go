package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

const (
	cStmt  = "delete from pockets where id = $1 RETURNING id;"
	cqStmt = "select * from pockets where id = $1"
)

func (h handler) Delete(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	id := c.Param("id")
	var p Pocket
	err := c.Bind(&p)

	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	stmt, err := h.db.Prepare(cqStmt)
	if err != nil {
		logger.Error("can't prepare query statement", zap.Error(err))
	}

	rows := stmt.QueryRow(id)
	err = rows.Scan(&p.ID, &p.Name, &p.Category, &p.Currency, &p.Balance)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Cloud pocket not found", err.Error())
	}

	if p.Balance > 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"Unable to delete this Cloud Pocket\n there is amount left in this Cloud Pocket, please move money out and try again", nil)
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, id).Scan(&lastInsertId)
	if err != nil {
		logger.Error("delete row error", zap.Error(err))
		return err
	}

	logger.Info("delete successfully", zap.Int64("id", lastInsertId))
	return c.JSON(http.StatusAccepted, "Cloud pocket deleted successfully")
}
