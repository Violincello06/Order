package handler

import (
	"github.com/labstack/echo/v4"
)

type SandboxHandler struct{}

func NewSandboxHandler() *SandboxHandler {
	return &SandboxHandler{}
}

func (h *SandboxHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/sandbox/pay/:payment_id", h.RenderSandboxPayment)
}

func (h *SandboxHandler) RenderSandboxPayment(c echo.Context) error {
	return c.File("public/sandbox.html")
}
