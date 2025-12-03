package com.payments.service;

import com.payments.dto.CreatePaymentRequest;
import com.payments.dto.PaymentResponse;
import com.payments.model.Payment;
import com.payments.model.PaymentStatus;
import com.payments.repository.PaymentRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
public class PaymentService {
    
    @Autowired
    private PaymentRepository paymentRepository;
    
    @Transactional
    public PaymentResponse createPayment(CreatePaymentRequest request, String userId) {
        Payment payment = new Payment();
        payment.setUserId(userId);
        payment.setAmount(request.getAmount());
        payment.setCurrency(request.getCurrency());
        payment.setMethod(request.getMethod());
        payment.setDescription(request.getDescription());
        payment.setMerchantId(request.getMerchantId());
        payment.setStatus(PaymentStatus.PENDING);
        payment.setTransactionId("TXN-" + UUID.randomUUID().toString());
        
        // Simulate payment processing
        payment.setStatus(PaymentStatus.PROCESSING);
        
        Payment saved = paymentRepository.save(payment);
        
        // Simulate async processing (in real app, send to queue)
        processPaymentAsync(saved.getId());
        
        return mapToResponse(saved);
    }
    
    public PaymentResponse getPaymentById(String id) {
        Payment payment = paymentRepository.findById(id)
            .orElseThrow(() -> new RuntimeException("Payment not found"));
        return mapToResponse(payment);
    }
    
    public List<PaymentResponse> getUserPayments(String userId) {
        return paymentRepository.findByUserId(userId)
            .stream()
            .map(this::mapToResponse)
            .collect(Collectors.toList());
    }
    
    private void processPaymentAsync(String paymentId) {
        // Simulate payment processing (would be async in production)
        new Thread(() -> {
            try {
                Thread.sleep(2000); // Simulate processing delay
                Payment payment = paymentRepository.findById(paymentId).orElse(null);
                if (payment != null) {
                    payment.setStatus(PaymentStatus.COMPLETED);
                    paymentRepository.save(payment);
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();
    }
    
    private PaymentResponse mapToResponse(Payment payment) {
        PaymentResponse response = new PaymentResponse();
        response.setId(payment.getId());
        response.setUserId(payment.getUserId());
        response.setAmount(payment.getAmount());
        response.setCurrency(payment.getCurrency());
        response.setStatus(payment.getStatus());
        response.setMethod(payment.getMethod());
        response.setDescription(payment.getDescription());
        response.setMerchantId(payment.getMerchantId());
        response.setTransactionId(payment.getTransactionId());
        response.setCreatedAt(payment.getCreatedAt());
        response.setUpdatedAt(payment.getUpdatedAt());
        return response;
    }
}
