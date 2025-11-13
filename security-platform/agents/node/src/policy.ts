import { PolicyRule } from './types';

export class PolicyEngine {
  private mode: 'observe' | 'block';
  private rules: PolicyRule[];

  constructor(config: { mode?: 'observe' | 'block'; rules?: PolicyRule[] }) {
    this.mode = config.mode || 'observe';
    this.rules = config.rules || [];
  }

  check(action: string, context: Record<string, any>): { allowed: boolean; reason?: string } {
    // In observe mode, always allow but log
    if (this.mode === 'observe') {
      return { allowed: true };
    }

    // Check rules
    for (const rule of this.rules) {
      if (rule.action === action || rule.action === '*') {
        if (this.evaluateCondition(rule.condition, context)) {
          return {
            allowed: rule.effect === 'allow',
            reason: rule.id,
          };
        }
      }
    }

    // Default allow if no rules match
    return { allowed: true };
  }

  private evaluateCondition(condition: string, context: Record<string, any>): boolean {
    // Simple JSONPath-like evaluation
    // In real implementation, use a proper expression evaluator
    try {
      // For now, simple key-value matching
      const parts = condition.split('==');
      if (parts.length === 2) {
        const key = parts[0].trim();
        const value = parts[1].trim().replace(/['"]/g, '');
        return context[key] === value;
      }
      return false;
    } catch (e) {
      return false;
    }
  }
}
