package cloudpocket

import "github.com/labstack/echo/v4"

type Cloudpockets struct {
	ID       int     `json:"id"`
	NAME     string  `json:"name"`
	CATAGORY string  `json:"catagory"`
	CURRENCY string  `json:"currency"`
	BALANCE  float64 `json:"balance"`
}

func GetAllCloudPocket(c echo.Context) error {
	dumpData := []Cloudpockets{
		{
			ID:       12345,
			NAME:     "Travel Fund",
			CATAGORY: "Vacation",
			CURRENCY: "THB",
			BALANCE:  100,
		},
		{
			ID:       67890,
			NAME:     "Savings",
			CATAGORY: "Emergency Fund",
			CURRENCY: "THB",
			BALANCE:  200,
		},
	}

	return c.JSON(200, dumpData)
}

func GetCloudPocketByID(c echo.Context) error {

	return nil
}
