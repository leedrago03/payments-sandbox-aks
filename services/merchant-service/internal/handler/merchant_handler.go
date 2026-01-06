package handler

import (
    "merchant-service/internal/model"
    "merchant-service/internal/service"
    
    "github.com/gofiber/fiber/v2"
)

type MerchantHandler struct {
    service *service.MerchantService
}

func NewMerchantHandler(service *service.MerchantService) *MerchantHandler {
    return &MerchantHandler{service: service}
}

// CreateMerchant creates a new merchant account
func (h *MerchantHandler) CreateMerchant(c *fiber.Ctx) error {
    var req model.CreateMerchantRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    merchant, err := h.service.CreateMerchant(&req)
    if err != nil {
        if err.Error() == "merchant with this email already exists" {
            return c.Status(409).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(merchant)
}

// GetMerchant retrieves merchant by ID
func (h *MerchantHandler) GetMerchant(c *fiber.Ctx) error {
    id := c.Params("id")
    
    merchant, err := h.service.GetMerchant(id)
    if err != nil {
        if err.Error() == "merchant not found" {
            return c.Status(404).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(merchant)
}

// UpdateMerchant updates merchant details
func (h *MerchantHandler) UpdateMerchant(c *fiber.Ctx) error {
    id := c.Params("id")
    
    var req model.CreateMerchantRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    merchant, err := h.service.UpdateMerchant(id, &req)
    if err != nil {
        if err.Error() == "merchant not found" {
            return c.Status(404).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(merchant)
}

// CreateAPIKey generates a new API key for merchant
func (h *MerchantHandler) CreateAPIKey(c *fiber.Ctx) error {
    merchantID := c.Params("id")
    
    var req model.CreateAPIKeyRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    apiKey, err := h.service.CreateAPIKey(merchantID, &req)
    if err != nil {
        if err.Error() == "merchant not found" {
            return c.Status(404).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(apiKey)
}

// GetAPIKeys retrieves all API keys for merchant
func (h *MerchantHandler) GetAPIKeys(c *fiber.Ctx) error {
    merchantID := c.Params("id")
    
    keys, err := h.service.GetAPIKeys(merchantID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(keys)
}

type VerifyAPIKeyRequest struct {
    APIKey string `json:"api_key"`
}

// VerifyAPIKey checks if an API key is valid
func (h *MerchantHandler) VerifyAPIKey(c *fiber.Ctx) error {
    var req VerifyAPIKeyRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    key, err := h.service.VerifyAPIKey(req.APIKey)
    if err != nil {
        if err.Error() == "invalid api key" {
            return c.Status(401).JSON(fiber.Map{"valid": false, "error": "Invalid API key"})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "valid":       true,
        "merchant_id": key.MerchantID,
    })
}
