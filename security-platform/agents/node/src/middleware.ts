import { Request, Response, NextFunction } from 'express';
import { Agent } from './agent';
import { MiddlewareOptions } from './types';

export function createMiddleware(agent: Agent, options: MiddlewareOptions = {}) {
  return (req: Request, res: Response, next: NextFunction) => {
    const startTime = Date.now();

    // Capture request metadata
    const requestMetadata = {
      method: req.method,
      path: req.path,
      headers: agent.redact(req.headers),
      query: agent.redact(req.query),
      body: agent.redact(req.body),
      ip: req.ip,
      userAgent: req.get('user-agent'),
    };

    // Check policy
    const policyResult = agent.checkPolicy('http_request', {
      method: req.method,
      path: req.path,
      ip: req.ip,
    });

    if (!policyResult.allowed && options.mode === 'block') {
      res.status(403).json({ error: 'Request blocked by security policy', reason: policyResult.reason });
      return;
    }

    // Capture response
    const originalSend = res.send;
    res.send = function (body: any) {
      const responseMetadata = {
        statusCode: res.statusCode,
        headers: agent.redact(res.getHeaders()),
        body: options.captureResponseBody ? agent.redact(body) : undefined,
        duration: Date.now() - startTime,
      };

      // Send to control plane (async, non-blocking)
      if (options.sendToControlPlane !== false) {
        sendEvent('http_request', {
          request: requestMetadata,
          response: responseMetadata,
          policy: policyResult,
        }, agent).catch(err => {
          console.error('Failed to send event:', err);
        });
      }

      return originalSend.call(this, body);
    };

    next();
  };
}

async function sendEvent(type: string, data: any, agent: Agent): Promise<void> {
  // In real implementation, send to control plane via HTTP or OTLP
  // For now, this is a placeholder
}
