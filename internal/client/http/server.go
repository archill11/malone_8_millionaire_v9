package http

import (
	"myapp/internal/service/tg_service"
	"myapp/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type (
	SerConfig struct {
		Port string
	}

	APIServer struct {
		Server *fiber.App
		l      *logger.Logger
		tgs    *tg_service.TgService
	}
)

func New(conf SerConfig, tgs *tg_service.TgService, l *logger.Logger) (*APIServer, error) {
	app := fiber.New()
	ser := &APIServer{
		Server: app,
		l:      l,
		tgs:    tgs,
	}

	return ser, nil
}