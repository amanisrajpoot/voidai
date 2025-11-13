/**
 * Security Platform Node.js Agent - Express Quick Start
 */

const express = require('express');
const { Agent, createMiddleware } = require('@security-platform/node-agent');

// Initialize agent
const agent = new Agent({
  controlPlaneUrl: process.env.SECURITY_PLATFORM_CONTROL_PLANE_URL || 'https://api.securityplatform.com',
  authKey: process.env.SECURITY_PLATFORM_AUTH_KEY,
  serviceName: 'my-api-service',
  environment: process.env.NODE_ENV || 'production',
  otlpEndpoint: process.env.SECURITY_PLATFORM_OTLP_ENDPOINT,
  policy: {
    mode: 'observe', // Start in observe mode, switch to 'block' after 48h
    autoEnableBlocking: false,
  },
  redactionRules: [
    { pattern: 'password|passwd', field: '*' },
    { pattern: '\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}', field: '*' }, // Credit card
  ],
});

// Start agent
agent.start();

const app = express();

// Add security platform middleware
app.use(createMiddleware(agent, {
  mode: 'observe', // or 'block' to enforce policies
  captureResponseBody: false, // Set to true if you need response body in events
}));

// Your routes
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

app.post('/api/login', express.json(), (req, res) => {
  // Request body will be automatically redacted before sending to control plane
  const { username, password } = req.body;
  
  // Check policy before processing
  const policyResult = agent.checkPolicy('login', {
    username,
    ip: req.ip,
  });
  
  if (!policyResult.allowed) {
    return res.status(403).json({ error: 'Login blocked by security policy' });
  }
  
  // Your login logic here
  res.json({ success: true });
});

app.listen(3000, () => {
  console.log('Server running on http://localhost:3000');
});

// Graceful shutdown
process.on('SIGTERM', () => {
  agent.stop();
  process.exit(0);
});
