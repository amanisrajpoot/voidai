/**
 * Sessions API routes
 */

import { Router, Request, Response } from 'express';
import { AuthRequest } from '../middleware/auth';

export const sessionsRouter = Router();

// In-memory store (replace with S3/MinIO in production)
const sessionsStore: Map<string, any> = new Map();

sessionsRouter.get('/:sessionId', (req: Request, res: Response) => {
  const { sessionId } = req.params;
  const session = sessionsStore.get(sessionId);

  if (!session) {
    return res.status(404).json({ error: 'Session not found' });
  }

  res.json(session);
});

sessionsRouter.get('/', (req: Request, res: Response) => {
  const { limit = 50, offset = 0 } = req.query;
  const sessions = Array.from(sessionsStore.values());

  const paginated = sessions.slice(Number(offset), Number(offset) + Number(limit));

  res.json({
    sessions: paginated,
    total: sessions.length,
    limit: Number(limit),
    offset: Number(offset),
  });
});

sessionsRouter.post('/:sessionId/replay', (req: Request, res: Response) => {
  const { sessionId } = req.params;
  const session = sessionsStore.get(sessionId);

  if (!session) {
    return res.status(404).json({ error: 'Session not found' });
  }

  // Return session replay data (rrweb events)
  res.json({
    sessionId,
    events: session.events || [],
    metadata: session.metadata || {},
  });
});
