/**
 * Authentication middleware
 */

import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';

const JWT_SECRET = process.env.JWT_SECRET || 'change-me-in-production';

export interface AuthRequest extends Request {
  auth?: {
    agentId?: string;
    serviceName?: string;
    permissions?: string[];
  };
}

export function authMiddleware(req: AuthRequest, res: Response, next: NextFunction) {
  const authHeader = req.headers.authorization;

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'Unauthorized: Missing or invalid authorization header' });
  }

  const token = authHeader.substring(7);

  try {
    const decoded = jwt.verify(token, JWT_SECRET) as any;
    req.auth = {
      agentId: decoded.agentId,
      serviceName: decoded.serviceName,
      permissions: decoded.permissions || [],
    };
    next();
  } catch (error) {
    return res.status(401).json({ error: 'Unauthorized: Invalid token' });
  }
}

export function generateToken(payload: { agentId: string; serviceName: string; permissions?: string[] }): string {
  return jwt.sign(payload, JWT_SECRET, { expiresIn: '7d' });
}
