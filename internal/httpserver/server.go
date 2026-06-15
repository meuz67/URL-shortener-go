package httpserver

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"url-shortener/internal/config"
	"url-shortener/internal/httpserver/handlers"
	"url-shortener/internal/storage"
)

type Server struct {
	cfg    config.Config
	engine *echo.Echo
}

func New(cfg config.Config, store storage.URLStore) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	urlHandler := handlers.NewURLHandler(cfg.BaseURL, store)
	RegisterRoutes(e, urlHandler)

	return &Server{
		cfg:    cfg,
		engine: e,
	}
}

func (s *Server) Start() error {
	if err := s.engine.Start(s.cfg.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
