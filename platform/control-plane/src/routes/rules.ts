/**
 * Rules API routes
 */

import { Router, Request, Response } from 'express';
import { AuthRequest } from '../middleware/auth';

export const rulesRouter = Router();

// In-memory store (replace with database in production)
const rulesStore: any[] = [
  {
    id: 'rule-1',
    name: 'SQL Injection Detection',
    pattern: '.*(union|select|drop|delete).*',
    action: 'alert',
    enabled: true,
  },
  {
    id: 'rule-2',
    name: 'Path Traversal Detection',
    pattern: '.*\\.\\./.*',
    action: 'block',
    enabled: true,
  },
];

rulesRouter.get('/', (req: Request, res: Response) => {
  res.json({
    rules: rulesStore.filter(r => r.enabled),
  });
});

rulesRouter.post('/', (req: Request, res: Response) => {
  const rule = {
    id: `rule-${Date.now()}`,
    ...req.body,
    createdAt: new Date().toISOString(),
  };

  rulesStore.push(rule);
  res.status(201).json(rule);
});

rulesRouter.put('/:ruleId', (req: Request, res: Response) => {
  const { ruleId } = req.params;
  const index = rulesStore.findIndex(r => r.id === ruleId);

  if (index === -1) {
    return res.status(404).json({ error: 'Rule not found' });
  }

  rulesStore[index] = { ...rulesStore[index], ...req.body, updatedAt: new Date().toISOString() };
  res.json(rulesStore[index]);
});

rulesRouter.delete('/:ruleId', (req: Request, res: Response) => {
  const { ruleId } = req.params;
  const index = rulesStore.findIndex(r => r.id === ruleId);

  if (index === -1) {
    return res.status(404).json({ error: 'Rule not found' });
  }

  rulesStore.splice(index, 1);
  res.status(204).send();
});
