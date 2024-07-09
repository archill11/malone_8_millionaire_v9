package main

import (
	"log"
	"myapp/config"
	api "myapp/internal/client/http"
	pg "myapp/internal/repository/pg"
	tg_service "myapp/internal/service/tg_service"
	"myapp/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

type application struct {
	config *config.Config
	server *api.APIServer
	logger *logger.Logger
	db     *pg.Database
	tgs    *tg_service.TgService
}

func main() {
	var err error
	app := &application{}

	app.logger = logger.New()
	app.config = config.Get()

	app.db, err = pg.New(app.config.Db, app.logger) // БД
	if err != nil {
		log.Fatal(err)
	}
	defer logFnError(app.db.CloseDb)

	app.tgs, err = tg_service.New(app.config.Tg, app.db, app.logger) // Tg Service
	if err != nil {
		log.Fatal(err)
	}

	app.server, err = api.New(app.config.Server, app.tgs, app.logger) // api server
	if err != nil {
		log.Fatal(err)
	}
	app.logger.Info("===============Listenning Server===============")
	go log.Fatal(app.server.Server.Listen(":" + app.config.Server.Port))

	defer func() {
		if err := app.server.Server.Shutdown(); err != nil {
			app.logger.Error(err)
		}
	}()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
	app.logger.Info("===============Server stopped===============")
}

func logFnError(fn func() error) {
	if err := fn(); err != nil {
		log.Println(err)
	}
}
