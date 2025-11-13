/**
 * Control Plane API
 * Main entry point for the security observability platform control plane
 */

import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import rateLimit from 'express-rate-limit';
import { eventsRouter } from './routes/events';
import { tracesRouter } from './routes/traces';
import { sessionsRouter } from './routes/sessions';
import { rulesRouter } from './routes/rules';
import { actionsRouter } from './routes/actions';
import { authMiddleware } from './middleware/auth';

const app = express();
const PORT = process.env.PORT || 3000;

// Security middleware
app.use(helmet());
app.use(cors({
  origin: process.env.ALLOWED_ORIGINS?.split(',') || '*',
  credentials: true,
}));

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 1000, // limit each IP to 1000 requests per windowMs
});
app.use('/api/', limiter);

// Body parsing
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// API routes
app.use('/api/v1/events', authMiddleware, eventsRouter);
app.use('/api/v1/traces', authMiddleware, tracesRouter);
app.use('/api/v1/sessions', authMiddleware, sessionsRouter);
app.use('/api/v1/rules', authMiddleware, rulesRouter);
app.use('/api/v1/actions', authMiddleware, actionsRouter);

// OTLP endpoints (for OpenTelemetry)
app.use('/v1/traces', authMiddleware, tracesRouter);
app.use('/v1/events', authMiddleware, eventsRouter);

// Error handling
app.use((err: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
  console.error('Error:', err);
  res.status(500).json({
    error: 'Internal server error',
    message: process.env.NODE_ENV === 'development' ? err.message : undefined,
  });
});

// Start server
app.listen(PORT, () => {
  console.log(`Control plane API running on port ${PORT}`);
  console.log(`Environment: ${process.env.NODE_ENV || 'development'}`);
});
