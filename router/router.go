package router

import (
	"database/sql"
	"github.com/kkgo-software-engineering/workshop/pocket"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/account"
	"github.com/kkgo-software-engineering/workshop/cloudpocket"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/featflag"
	"github.com/kkgo-software-engineering/workshop/healthchk"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Use(mlog.Middleware(logger))
	//e.Use(middleware.BasicAuth(mw.Authenicate()))

	hHealthChk := healthchk.New(db)
	e.GET("/healthz", hHealthChk.Check)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	hAccount := account.New(cfg.FeatureFlag, db)
	hPocket := pocket.New(cfg.FeatureFlag, db)
	e.POST("/accounts", hAccount.Create)
	e.POST("/cloud-pockets/transfer", hPocket.Transfer)

	cloudpocketDelete := cloudpocket.New(cfg.FeatureFlag, db)
	e.DELETE("/cloud-pocket/:id", cloudpocketDelete.Delete)

	hFeatFlag := featflag.New(cfg)
	e.GET("/features", hFeatFlag.List)

	return e
}
