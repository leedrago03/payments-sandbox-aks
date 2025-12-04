package service

import (
    "audit-service/internal/model"
    "audit-service/internal/repository"
    "time"
    
    "github.com/google/uuid"
)

type AuditService struct {
    repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
    return &AuditService{repo: repo}
}

func (s *AuditService) LogEvent(req *model.CreateAuditLogRequest) (*model.AuditLog, error) {
    log := &model.AuditLog{
        ID:         uuid.New().String(),
        EventType:  req.EventType,
        EntityType: req.EntityType,
        EntityID:   req.EntityID,
        ActorType:  req.ActorType,
        ActorID:    req.ActorID,
        Action:     req.Action,
        Details:    req.Details,
        IPAddress:  req.IPAddress,
        UserAgent:  req.UserAgent,
        Success:    req.Success,
        ErrorMsg:   req.ErrorMsg,
        CreatedAt:  time.Now(),
    }
    
    if err := s.repo.Create(log); err != nil {
        return nil, err
    }
    
    return log, nil
}

func (s *AuditService) QueryLogs(params *model.AuditLogQueryParams) ([]model.AuditLog, error) {
    return s.repo.Query(params)
}

func (s *AuditService) GetEntityLogs(entityID string, limit int) ([]model.AuditLog, error) {
    if limit <= 0 {
        limit = 50
    }
    return s.repo.GetByEntityID(entityID, limit)
}

func (s *AuditService) GetStats() (map[string]int, error) {
    return s.repo.GetStats()
}
