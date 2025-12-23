package handler

import (
	"payment-service/internal/model"
	"payment-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: svc,
	}
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

	payment, err := h.service.Authorize(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
