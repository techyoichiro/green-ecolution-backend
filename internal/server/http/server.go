package http

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

type HTTPError struct {
	Error  string `json:"error"`
	Code   int    `json:"code"`
	Path   string `json:"path"`
	Method string `json:"method"`
} // @Name HTTPError

type Server struct {
	cfg      *config.Config
	services *service.Services
}

func NewServer(cfg *config.Config, services *service.Services) *Server {
	return &Server{
		cfg:      cfg,
		services: services,
	}
}

func (s *Server) Run(ctx context.Context) error {
	app := fiber.New(fiber.Config{
		AppName:      s.cfg.Dashboard.Title,
		ServerHeader: s.cfg.Dashboard.Title,
		ErrorHandler: errorHandler,
	})

	app.Mount("/", s.middleware(s.publicRoutes, s.privateRoutes))

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down HTTP Server")
		if err := app.Shutdown(); err != nil {
			fmt.Println("Error while shutting down HTTP Server:", err)
		}
	}()

	return app.Listen(fmt.Sprintf(":%d", s.cfg.Server.Port))
}

func errorHandler(c *fiber.Ctx, err error) error {
	c.Status(fiber.StatusInternalServerError)
	var e *fiber.Error
	if errors.As(err, &e) {
		return c.JSON(HTTPError{
			e.Message,
			e.Code,
			c.Path(),
			c.Method(),
		})
	}
	return nil
}
