package com.payments.dto;

import com.payments.model.PaymentMethod;
import com.payments.model.PaymentStatus;
import lombok.Data;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
public class PaymentResponse {
    private String id;
    private String userId;
    private BigDecimal amount;
    private String currency;
    private PaymentStatus status;
    private PaymentMethod method;
    private String description;
    private String merchantId;
    private String transactionId;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}
