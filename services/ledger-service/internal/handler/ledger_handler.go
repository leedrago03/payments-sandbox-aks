package handler

import (
    "ledger-service/internal/model"
    "ledger-service/internal/service"
    
    "github.com/gofiber/fiber/v2"
)

type LedgerHandler struct {
    service *service.LedgerService
}

func NewLedgerHandler(service *service.LedgerService) *LedgerHandler {
    return &LedgerHandler{service: service}
}

func (h *LedgerHandler) CreateEntries(c *fiber.Ctx) error {
    var req model.CreateEntriesRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    if len(req.Entries) < 2 {
        return c.Status(400).JSON(fiber.Map{"error": "At least 2 entries required for double-entry"})
    }
    
    response, err := h.service.CreateEntries(&req)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(response)
}

func (h *LedgerHandler) GetPaymentEntries(c *fiber.Ctx) error {
    paymentID := c.Params("paymentId")
    
    entries, err := h.service.GetPaymentEntries(paymentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(entries)
}

func (h *LedgerHandler) GetAccountBalance(c *fiber.Ctx) error {
    accountID := c.Params("accountId")
    currency := c.Query("currency", "USD")
    
    balance, err := h.service.GetAccountBalance(accountID, currency)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(balance)
}

func (h *LedgerHandler) GetAllBalances(c *fiber.Ctx) error {
    balances, err := h.service.GetAllBalances()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(balances)
}
