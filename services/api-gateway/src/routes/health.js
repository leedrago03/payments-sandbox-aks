const express = require('express');
const router = express.Router();

// Kubernetes liveness probe
router.get('/liveness', (req, res) => {
  res.status(200).json({ status: 'alive' });
});

// Kubernetes readiness probe
router.get('/readiness', (req, res) => {
  // Check dependencies (Redis, downstream services)
  // For now, simple response
  res.status(200).json({ 
    status: 'ready',
    service: 'api-gateway',
    timestamp: new Date().toISOString()
  });
});

module.exports = router;
