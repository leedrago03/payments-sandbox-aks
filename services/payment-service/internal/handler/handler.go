package handler

import (
	"payment-service/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type PaymentHandler struct {
	// In a real scenario, we'd inject the service here
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/v1")
	v1.Post("/payments", h.AuthorizePayment)
	v1.Post("/payments/:id/capture", h.CapturePayment)
	v1.Get("/payments/:id", h.GetPayment)
}

func (h *PaymentHandler) AuthorizePayment(c *fiber.Ctx) error {
	var req model.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Mock logic for now
	payment := model.Payment{
		ID:             uuid.New().String(),
		MerchantID:     req.MerchantID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Token:          req.Token,
		Status:         model.StatusAuthorized,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}

func (h *PaymentHandler) CapturePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	// Mock logic
	return c.JSON(fiber.Map{
		"id":     id,
		"status": model.StatusCaptured,
	})
}

func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	id := c.Params("id")
	// Mock logic
	return c.JSON(fiber.Map{
		"id":     id,
		"status": model.StatusAuthorized,
	})
}
