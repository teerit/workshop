package pocket

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *handler) GetAllCloudPocket(c echo.Context) error {
	logger := mlog.L(c)
	stmt, err := h.db.Prepare("SELECT id, name, category, currency, balance FROM pockets")
	defer logger.Sync()
	if err != nil {
		logger.Error("Can't prepare query all cloud pockets statment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, errorResp{ErrorMessage: "Can't prepare query all cloud pockets statment"})
	}

	rows, err := stmt.Query()
	if err != nil {
		logger.Error("Can't query all cloud pockets", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, errorResp{ErrorMessage: "Can't query all cloud pockets"})
	}

	pockets := []pocket{}
	for rows.Next() {
		var cpk pocket
		err = rows.Scan(&cpk.ID, &cpk.Name, &cpk.Category, &cpk.Currency, &cpk.Balance)
		if err != nil {
			logger.Error("Can't scan cloud pocket", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, errorResp{ErrorMessage: "Can't scan cloud pocket"})
		}
		pockets = append(pockets, cpk)
	}
	return c.JSON(http.StatusOK, pockets)
}

func (h *handler) GetCloudPocketByID(c echo.Context) error {
	logger := mlog.L(c)
	id := c.Param("id")
	stmt, err := h.db.Prepare("SELECT id, name, category, currency, balance FROM pockets WHERE id = $1")
	defer logger.Sync()
	if err != nil {
		logger.Error("Can't prepare query cloud pocket statment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, errorResp{ErrorMessage: "Can't prepare query cloud pocket statment"})
	}
	row := stmt.QueryRow(id)
	cpk := pocket{}
	err = row.Scan(&cpk.ID, &cpk.Name, &cpk.Category, &cpk.Currency, &cpk.Balance)
	switch err {
	case sql.ErrNoRows:
		logger.Error("Cloud Pocket Data Not Found", zap.Error(err))
		return c.JSON(http.StatusNoContent, errorResp{ErrorMessage: "Cloud Pocket Data Not Found"})
	case nil:
		return c.JSON(http.StatusOK, cpk)
	default:
		logger.Error("Can't scan cloud pocket", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, errorResp{ErrorMessage: "Can't scan cloud pocket"})
	}
}
