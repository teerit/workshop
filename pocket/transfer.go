package pocket

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type transferDto struct {
	SourcePocketId string  `json:"source_cloud_pocket_id"`
	DestPocketId   string  `json:"destination_cloud_pocket_id"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description"`
}

type transferResponse struct {
	TransactionId          int     `json:"transaction_id"`
	SourceCloudPocket      *Pocket `json:"source_cloud_pocket"`
	DestinationCloudPocket *Pocket `json:"destination_cloud_pocket"`
	Status                 string  `json:"status"`
}

func (h *handler) Transfer(c echo.Context) error {
	// init db and variable
	tDto := &transferDto{}
	sourcePocket := &Pocket{}
	destPocket := &Pocket{}
	db, err := h.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}
	defer func() {
		if err != nil {
			_ = db.Rollback()
		}
		_ = db.Commit()
	}()
	tfService := newTransferService(db)

	// bind dto
	err = c.Bind(tDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp{
			Status:       "Failed",
			ErrorMessage: "Bad request",
		})
	}

	// query and bind source pocket
	err = tfService.findPocket(tDto.SourcePocketId, sourcePocket)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not found source pocket",
		})
	}
	// query and bind destination pocket
	err = tfService.findPocket(tDto.DestPocketId, destPocket)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not found destination pocket",
		})
	}

	isBalance, err := tfService.balanceCheck(tDto, sourcePocket.Balance)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}
	if isBalance {
		return c.JSON(http.StatusBadRequest, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not enough balance in the source cloud pocket",
		})
	}

	// update amount source and destination
	err = tfService.updatePocket(sourcePocket, sourcePocket.Balance-tDto.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	err = tfService.updatePocket(destPocket, destPocket.Balance+tDto.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	tId, err := tfService.insertTransaction(tDto, "Success")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	resp := &transferResponse{
		TransactionId:          tId,
		SourceCloudPocket:      sourcePocket,
		DestinationCloudPocket: destPocket,
		Status:                 "Success",
	}

	return c.JSON(http.StatusOK, resp)
}
