package pocket

import "github.com/labstack/echo/v4"

func GetAllCloudPocket(c echo.Context) error {
	dumpData := []pocket{
		{
			ID:       12345,
			Name:     "Travel Fund",
			Category: "Vacation",
			Currency: "THB",
			Balance:  100,
		},
		{
			ID:       67890,
			Name:     "Savings",
			Category: "Emergency Fund",
			Currency: "THB",
			Balance:  200,
		},
	}

	return c.JSON(200, dumpData)
}

func GetCloudPocketByID(c echo.Context) error {

	return nil
}
