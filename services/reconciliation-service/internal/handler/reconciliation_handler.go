package handler

import (
    "reconciliation-service/internal/model"
    "reconciliation-service/internal/service"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
)

type ReconciliationHandler struct {
    service *service.ReconciliationService
}

func NewReconciliationHandler(service *service.ReconciliationService) *ReconciliationHandler {
    return &ReconciliationHandler{service: service}
}

func (h *ReconciliationHandler) CreateReconciliation(c *fiber.Ctx) error {
    var req model.CreateReportRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    summary, err := h.service.CreateReconciliation(&req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(summary)
}

func (h *ReconciliationHandler) GetReport(c *fiber.Ctx) error {
    reportID := c.Params("reportId")
    
    summary, err := h.service.GetReport(reportID)
    if err != nil {
        if err.Error() == "report not found" {
            return c.Status(404).JSON(fiber.Map{"error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(summary)
}

func (h *ReconciliationHandler) GetAllReports(c *fiber.Ctx) error {
    limit := 50
    if l := c.Query("limit"); l != "" {
        if parsed, err := strconv.Atoi(l); err == nil {
            limit = parsed
        }
    }
    
    reports, err := h.service.GetAllReports(limit)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(reports)
}
