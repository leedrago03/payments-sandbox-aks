package repository

import (
    "audit-service/internal/model"
    "database/sql"
    "fmt"
)

type AuditRepository struct {
    db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
    return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log *model.AuditLog) error {
    query := `
        INSERT INTO audit_logs (id, event_type, entity_type, entity_id, actor_type, 
                               actor_id, action, details, ip_address, user_agent, 
                               success, error_msg, signature, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `
    
    _, err := r.db.Exec(query, log.ID, log.EventType, log.EntityType, log.EntityID,
        log.ActorType, log.ActorID, log.Action, log.Details, log.IPAddress,
        log.UserAgent, log.Success, log.ErrorMsg, log.Signature, log.CreatedAt)
    
    return err
}

func (r *AuditRepository) Query(params *model.AuditLogQueryParams) ([]model.AuditLog, error) {
    query := `
        SELECT id, event_type, entity_type, entity_id, actor_type, 
               actor_id, action, details, ip_address, user_agent, 
               success, error_msg, signature, created_at 
        FROM audit_logs WHERE 1=1`
    args := []interface{}{}
    argCount := 1
    
    // ... existing filtering logic ...
    if params.EntityType != "" {
        query += fmt.Sprintf(" AND entity_type = $%d", argCount)
        args = append(args, params.EntityType)
        argCount++
    }
    
    if params.EntityID != "" {
        query += fmt.Sprintf(" AND entity_id = $%d", argCount)
        args = append(args, params.EntityID)
        argCount++
    }
    
    if params.ActorID != "" {
        query += fmt.Sprintf(" AND actor_id = $%d", argCount)
        args = append(args, params.ActorID)
        argCount++
    }
    
    if params.EventType != "" {
        query += fmt.Sprintf(" AND event_type = $%d", argCount)
        args = append(args, params.EventType)
        argCount++
    }
    
    if !params.StartDate.IsZero() {
        query += fmt.Sprintf(" AND created_at >= $%d", argCount)
        args = append(args, params.StartDate)
        argCount++
    }
    
    if !params.EndDate.IsZero() {
        query += fmt.Sprintf(" AND created_at <= $%d", argCount)
        args = append(args, params.EndDate)
        argCount++
    }
    
    query += " ORDER BY created_at DESC"
    
    if params.Limit > 0 {
        query += fmt.Sprintf(" LIMIT $%d", argCount)
        args = append(args, params.Limit)
    } else {
        query += " LIMIT 100"
    }
    
    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var logs []model.AuditLog
    for rows.Next() {
        var log model.AuditLog
        err := rows.Scan(&log.ID, &log.EventType, &log.EntityType, &log.EntityID,
            &log.ActorType, &log.ActorID, &log.Action, &log.Details,
            &log.IPAddress, &log.UserAgent, &log.Success, &log.ErrorMsg, 
            &log.Signature, &log.CreatedAt)
        if err != nil {
            return nil, err
        }
        logs = append(logs, log)
    }
    
    return logs, nil
}

func (r *AuditRepository) GetByEntityID(entityID string, limit int) ([]model.AuditLog, error) {
    query := `
        SELECT id, event_type, entity_type, entity_id, actor_type, 
               actor_id, action, details, ip_address, user_agent, 
               success, error_msg, signature, created_at 
        FROM audit_logs 
        WHERE entity_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2
    `
    
    rows, err := r.db.Query(query, entityID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var logs []model.AuditLog
    for rows.Next() {
        var log model.AuditLog
        err := rows.Scan(&log.ID, &log.EventType, &log.EntityType, &log.EntityID,
            &log.ActorType, &log.ActorID, &log.Action, &log.Details,
            &log.IPAddress, &log.UserAgent, &log.Success, &log.ErrorMsg, 
            &log.Signature, &log.CreatedAt)
        if err != nil {
            return nil, err
        }
        logs = append(logs, log)
    }
    
    return logs, nil
}


func (r *AuditRepository) GetStats() (map[string]int, error) {
    query := `
        SELECT event_type, COUNT(*) as count 
        FROM audit_logs 
        GROUP BY event_type
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    stats := make(map[string]int)
    for rows.Next() {
        var eventType string
        var count int
        if err := rows.Scan(&eventType, &count); err != nil {
            return nil, err
        }
        stats[eventType] = count
    }
    
    return stats, nil
}
