package handler

import (
    "tokenization-service/internal/model"
    "tokenization-service/internal/service"
    
    "github.com/gofiber/fiber/v2"
)

type TokenHandler struct {
    service *service.TokenService
}

func NewTokenHandler(service *service.TokenService) *TokenHandler {
    return &TokenHandler{service: service}
}

func (h *TokenHandler) Tokenize(c *fiber.Ctx) error {
    var req model.TokenizeRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    response, err := h.service.Tokenize(&req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(response)
}

func (h *TokenHandler) GetToken(c *fiber.Ctx) error {
    token := c.Params("token")
    
    result, err := h.service.GetToken(token)
    if err != nil {
        if err.Error() == "token not found" {
            return c.Status(404).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(result)
}

func (h *TokenHandler) Detokenize(c *fiber.Ctx) error {
    token := c.Params("token")
    
    pan, err := h.service.DetokenizePAN(token)
    if err != nil {
        if err.Error() == "token not found" {
            return c.Status(404).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(fiber.Map{"pan": pan})
}
