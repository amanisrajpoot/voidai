/**
 * Traces API routes (OTLP)
 */

import { Router, Request, Response } from 'express';
import { AuthRequest } from '../middleware/auth';

export const tracesRouter = Router();

// In-memory store (replace with OpenTelemetry collector in production)
const tracesStore: any[] = [];

tracesRouter.post('/', (req: AuthRequest, res: Response) => {
  // OTLP trace export endpoint
  // In production, this would forward to OpenTelemetry Collector
  const traces = req.body;

  tracesStore.push({
    ...traces,
    receivedAt: new Date().toISOString(),
    agentId: req.auth?.agentId,
    serviceName: req.auth?.serviceName,
  });

  res.status(200).send();
});

tracesRouter.get('/', (req: Request, res: Response) => {
  const { limit = 100, serviceName, traceId } = req.query;

  let filtered = tracesStore;

  if (serviceName) {
    filtered = filtered.filter((t: any) => t.serviceName === serviceName);
  }

  if (traceId) {
    filtered = filtered.filter((t: any) => t.resourceSpans?.some((rs: any) =>
      rs.instrumentationLibrarySpans?.some((ils: any) =>
        ils.spans?.some((s: any) => s.traceId === traceId)
      )
    ));
  }

  const paginated = filtered.slice(0, Number(limit));

  res.json({
    traces: paginated,
    total: filtered.length,
  });
});
