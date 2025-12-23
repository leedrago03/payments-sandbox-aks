package service

import (
    "audit-service/internal/model"
    "audit-service/internal/repository"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "github.com/payments-sandbox/pkg/crypto"
)

type AuditService struct {
    repo *repository.AuditRepository
    hmacKey []byte
}

func NewAuditService(repo *repository.AuditRepository, hmacKey string) *AuditService {
    return &AuditService{
        repo: repo,
        hmacKey: []byte(hmacKey),
    }
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
    
    // Calculate signature for immutability
    log.Signature = s.calculateSignature(log)
    
    if err := s.repo.Create(log); err != nil {
        return nil, err
    }
    
    return log, nil
}

func (s *AuditService) calculateSignature(log *model.AuditLog) string {
    // Concatenate key fields to create a message for signing
    // Using a stable format (CSV-like) for the signature base
    message := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%t|%s|%s",
        log.ID,
        log.EventType,
        log.EntityID,
        log.ActorID,
        log.Action,
        log.Details,
        log.Success,
        log.ErrorMsg,
        log.CreatedAt.Format(time.RFC3339),
    )
    
    return crypto.GenerateHMAC([]byte(message), s.hmacKey)
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
