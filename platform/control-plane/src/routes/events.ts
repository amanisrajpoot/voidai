/**
 * Events API routes
 */

import { Router, Request, Response } from 'express';
import { AuthRequest } from '../middleware/auth';

export const eventsRouter = Router();

// In-memory store (replace with database in production)
const eventsStore: any[] = [];

eventsRouter.post('/', (req: AuthRequest, res: Response) => {
  const { events, sessionId } = req.body;

  if (!Array.isArray(events)) {
    return res.status(400).json({ error: 'events must be an array' });
  }

  // Enrich events with metadata
  const enrichedEvents = events.map((event: any) => ({
    ...event,
    receivedAt: new Date().toISOString(),
    agentId: req.auth?.agentId,
    serviceName: req.auth?.serviceName,
    sessionId: sessionId || event.sessionId,
  }));

  eventsStore.push(...enrichedEvents);

  // TODO: Store in OpenSearch/ClickHouse
  // TODO: Trigger rule evaluation
  // TODO: Send alerts if needed

  res.json({
    success: true,
    received: enrichedEvents.length,
  });
});

eventsRouter.get('/', (req: Request, res: Response) => {
  const { limit = 100, offset = 0, type, sessionId } = req.query;

  let filtered = eventsStore;

  if (type) {
    filtered = filtered.filter((e: any) => e.type === type);
  }

  if (sessionId) {
    filtered = filtered.filter((e: any) => e.sessionId === sessionId);
  }

  const paginated = filtered.slice(Number(offset), Number(offset) + Number(limit));

  res.json({
    events: paginated,
    total: filtered.length,
    limit: Number(limit),
    offset: Number(offset),
  });
});

eventsRouter.get('/:eventId', (req: Request, res: Response) => {
  const { eventId } = req.params;
  const event = eventsStore.find((e: any) => e.id === eventId);

  if (!event) {
    return res.status(404).json({ error: 'Event not found' });
  }

  res.json(event);
});
