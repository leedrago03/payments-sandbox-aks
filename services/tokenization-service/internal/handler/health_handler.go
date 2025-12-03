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
    if err := h.db.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status":   "not ready",
            "database": "disconnected",
        })
    }
    
    return c.JSON(fiber.Map{
        "status":   "ready",
        "service":  "tokenization-service",
        "database": "connected",
    })
}
