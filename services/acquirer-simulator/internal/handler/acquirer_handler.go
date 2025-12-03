package handler

import (
    "acquirer-simulator/internal/model"
    "acquirer-simulator/internal/service"
    
    "github.com/gofiber/fiber/v2"
)

type AcquirerHandler struct {
    service *service.AcquirerService
}

func NewAcquirerHandler(service *service.AcquirerService) *AcquirerHandler {
    return &AcquirerHandler{service: service}
}

func (h *AcquirerHandler) Authorize(c *fiber.Ctx) error {
    var req model.AuthRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    response, err := h.service.Authorize(&req)
    if err != nil {
        return c.Status(504).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(response)
}

func (h *AcquirerHandler) Capture(c *fiber.Ctx) error {
    var req model.CaptureRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    response, err := h.service.Capture(&req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(response)
}

func (h *AcquirerHandler) Refund(c *fiber.Ctx) error {
    var req model.RefundRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    response, err := h.service.Refund(&req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(response)
}

func (h *AcquirerHandler) Health(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "healthy", "service": "acquirer-simulator"})
}
