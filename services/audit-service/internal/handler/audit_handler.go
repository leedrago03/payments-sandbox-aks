package handler

import (
    "audit-service/internal/model"
    "audit-service/internal/service"
    "strconv"
    "time"
    
    "github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
    service *service.AuditService
}

func NewAuditHandler(service *service.AuditService) *AuditHandler {
    return &AuditHandler{service: service}
}

func (h *AuditHandler) CreateLog(c *fiber.Ctx) error {
    var req model.CreateAuditLogRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    // Extract IP and User-Agent from request
    if req.IPAddress == "" {
        req.IPAddress = c.IP()
    }
    if req.UserAgent == "" {
        req.UserAgent = c.Get("User-Agent")
    }
    
    log, err := h.service.LogEvent(&req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(log)
}

func (h *AuditHandler) QueryLogs(c *fiber.Ctx) error {
    params := &model.AuditLogQueryParams{
        EntityType: c.Query("entity_type"),
        EntityID:   c.Query("entity_id"),
        ActorID:    c.Query("actor_id"),
        EventType:  model.EventType(c.Query("event_type")),
    }
    
    if startDate := c.Query("start_date"); startDate != "" {
        if t, err := time.Parse(time.RFC3339, startDate); err == nil {
            params.StartDate = t
        }
    }
    
    if endDate := c.Query("end_date"); endDate != "" {
        if t, err := time.Parse(time.RFC3339, endDate); err == nil {
            params.EndDate = t
        }
    }
    
    if limit := c.Query("limit"); limit != "" {
        if l, err := strconv.Atoi(limit); err == nil {
            params.Limit = l
        }
    }
    
    logs, err := h.service.QueryLogs(params)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(logs)
}

func (h *AuditHandler) GetEntityLogs(c *fiber.Ctx) error {
    entityID := c.Params("entityId")
    limit := 50
    
    if l := c.Query("limit"); l != "" {
        if parsed, err := strconv.Atoi(l); err == nil {
            limit = parsed
        }
    }
    
    logs, err := h.service.GetEntityLogs(entityID, limit)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(logs)
}

func (h *AuditHandler) GetStats(c *fiber.Ctx) error {
    stats, err := h.service.GetStats()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(stats)
}
