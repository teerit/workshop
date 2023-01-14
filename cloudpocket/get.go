package cloudpocket

import "github.com/labstack/echo/v4"

type Cloudpockets struct {
	ID       int     `json:"id"`
	NAME     string  `json:"name"`
	CATAGORY string  `json:"catagory`
	CURRENCY string  `json:"currency"`
	BALANCE  float64 `json:"balance"`
}

func GetAllCloudPocket(c echo.Context) error {
	dumpData := []Cloudpockets{
		{
			"id":       "12345",
			"name":     "Travel Fund",
			"category": "Vacation",
			"currency": "THB",
			"balance":  100,
		},
		{
			"id":       "67890",
			"name":     "Savings",
			"category": "Emergency Fund",
			"currency": "THB",
			"balance":  200,
		}
	}

	return c.JSON(200, dumpData)
}
