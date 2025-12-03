package handler

import (
    "database/sql"
    "github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
    db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
    return &HealthHandler{db: db}
}

func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "alive"})
}

func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
    // Check database connection
    if err := h.db.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "not ready",
            "database": "disconnected",
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "status": "ready",
        "service": "merchant-service",
        "database": "connected",
    })
}
