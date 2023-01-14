package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/router"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New().All()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg.DBConnection)

	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		logger.Fatal("unable to configure database", zap.Error(err))
	}

	err = initTable(sql)
	if err != nil {
		logger.Fatal("error init-db", zap.Error(err))
	}

	e := router.RegRoute(cfg, logger, sql)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port)

	go func() {
		err := e.Start(addr)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("unexpected shutdown the server", zap.Error(err))
		}
		logger.Info("gracefully shutdown the server")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gCtx := context.Background()
	ctx, cancel := context.WithTimeout(gCtx, 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("unexpected shutdown the server", zap.Error(err))
	}
}

func initTable(db *sql.DB) error {
	createTb := `
			CREATE TABLE IF NOT EXISTS pockets(
				id SERIAL PRIMARY KEY,
				name TEXT,
				category TEXT,
				currency TEXT,
				balance float8
			);
`
	createTransaction := `
		CREATE TABLE IF NOT EXISTS transactions(
			id SERIAL PRIMARY KEY,
			source_pid INT,
			dest_pid INT,
			amount float8,
			description TEXT,
			date timestamp,
			status TEXT
		);
`
	_, err := db.Exec(createTb)
	if err != nil {
		return err
	}

	_, err = db.Exec(createTransaction)
	if err != nil {
		return err
	}
	return nil
}
