const express = require('express');
const axios = require('axios');
const { authenticateToken } = require('../middleware/auth');
const router = express.Router();

const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://payment-service:8080';

// All payment routes require authentication
router.use(authenticateToken);

// Create payment
router.post('/', async (req, res) => {
  try {
    // Forward request to Payment Service
    const response = await axios.post(`${PAYMENT_SERVICE_URL}/api/payments`, {
      ...req.body,
      userId: req.user.userId
    });
    res.json(response.data);
  } catch (error) {
    res.status(error.response?.status || 500).json({
      error: error.response?.data || error.message
    });
  }
});

// Get payment by ID
router.get('/:id', async (req, res) => {
  try {
    const response = await axios.get(
      `${PAYMENT_SERVICE_URL}/api/payments/${req.params.id}`,
      { headers: { 'X-User-Id': req.user.userId } }
    );
    res.json(response.data);
  } catch (error) {
    res.status(error.response?.status || 500).json({
      error: error.response?.data || error.message
    });
  }
});

// List user payments
router.get('/', async (req, res) => {
  try {
    const response = await axios.get(
      `${PAYMENT_SERVICE_URL}/api/payments`,
      { 
        headers: { 'X-User-Id': req.user.userId },
        params: req.query
      }
    );
    res.json(response.data);
  } catch (error) {
    res.status(error.response?.status || 500).json({
      error: error.response?.data || error.message
    });
  }
});

module.exports = router;
