package server

import (
	"context"
	"encoding/json"
	config "rent-application/configs"
	"rent-application/shared/web"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

type Server struct {
	app    *fiber.App
	logger *zerolog.Logger
	cfg    config.Config
}

func New(cfg config.Config, logger *zerolog.Logger) (*Server, error) {
	app := fiber.New(fiber.Config{
		ErrorHandler:          web.ErrorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		BodyLimit:             15 * 1024 * 1024,
		DisableStartupMessage: true,
	})

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool { return true },
	}))

	// Semua middleware dan route diatur di SetupRouter
	SetupRouter(app, cfg, logger)

	return &Server{
		app:    app,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (s *Server) Start() error {
	return s.app.Listen(s.cfg.App.Host + ":" + s.cfg.App.Port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
