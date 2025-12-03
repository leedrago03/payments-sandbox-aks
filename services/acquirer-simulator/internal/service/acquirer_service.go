package service

import (
    "acquirer-simulator/internal/model"
    "errors"
    "math/rand"
    "strconv"
    "time"
    
    "github.com/google/uuid"
)

type AcquirerService struct {
    successRate int
    timeoutRate int
}

func NewAcquirerService(successRate, timeoutRate string) *AcquirerService {
    sr, _ := strconv.Atoi(successRate)
    tr, _ := strconv.Atoi(timeoutRate)
    
    rand.Seed(time.Now().UnixNano())
    
    return &AcquirerService{
        successRate: sr,
        timeoutRate: tr,
    }
}

func (s *AcquirerService) Authorize(req *model.AuthRequest) (*model.AuthResponse, error) {
    // Simulate processing delay
    time.Sleep(time.Millisecond * time.Duration(100+rand.Intn(200)))
    
    transactionID := "TXN_" + uuid.New().String()
    
    // Simulate random outcomes
    outcome := rand.Intn(100)
    
    if outcome < s.timeoutRate {
        return nil, errors.New("gateway timeout")
    } else if outcome < s.successRate+s.timeoutRate {
        // Success
        return &model.AuthResponse{
            Status:        "APPROVED",
            TransactionID: transactionID,
            AuthCode:      generateAuthCode(),
            Amount:        req.Amount,
            Currency:      req.Currency,
        }, nil
    } else {
        // Declined
        return &model.AuthResponse{
            Status:        "DECLINED",
            TransactionID: transactionID,
            DeclineReason: getRandomDeclineReason(),
            Amount:        req.Amount,
            Currency:      req.Currency,
        }, nil
    }
}

func (s *AcquirerService) Capture(req *model.CaptureRequest) (*model.CaptureResponse, error) {
    // Simulate processing
    time.Sleep(time.Millisecond * time.Duration(50+rand.Intn(100)))
    
    // 95% success rate for captures
    if rand.Intn(100) < 95 {
        return &model.CaptureResponse{
            Status:         "CAPTURED",
            TransactionID:  req.TransactionID,
            CapturedAmount: req.Amount,
        }, nil
    }
    
    return &model.CaptureResponse{
        Status:        "FAILED",
        TransactionID: req.TransactionID,
    }, nil
}

func (s *AcquirerService) Refund(req *model.RefundRequest) (*model.RefundResponse, error) {
    // Simulate processing
    time.Sleep(time.Millisecond * time.Duration(100+rand.Intn(150)))
    
    // 98% success rate for refunds
    if rand.Intn(100) < 98 {
        return &model.RefundResponse{
            Status:        "REFUNDED",
            TransactionID: req.TransactionID,
            RefundID:      "REF_" + uuid.New().String(),
            RefundAmount:  req.Amount,
        }, nil
    }
    
    return &model.RefundResponse{
        Status:        "FAILED",
        TransactionID: req.TransactionID,
    }, nil
}

func generateAuthCode() string {
    return strconv.Itoa(100000 + rand.Intn(900000))
}

func getRandomDeclineReason() string {
    reasons := []string{
        "INSUFFICIENT_FUNDS",
        "EXPIRED_CARD",
        "INVALID_CVV",
        "DO_NOT_HONOR",
        "SUSPECTED_FRAUD",
        "LIMIT_EXCEEDED",
    }
    return reasons[rand.Intn(len(reasons))]
}
