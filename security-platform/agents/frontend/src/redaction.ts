import { RedactionRule } from './types';

const DEFAULT_BLOCKLIST = [
  'password',
  'passwd',
  'pwd',
  'secret',
  'token',
  'api_key',
  'apikey',
  'auth',
  'credit_card',
  'card_number',
  'cvv',
  'ssn',
  'social_security',
  'pin',
  'passcode',
];

export class Redactor {
  private rules: RedactionRule[];
  private patterns: RegExp[];

  constructor(rules: RedactionRule[] = []) {
    this.rules = rules;
    this.patterns = DEFAULT_BLOCKLIST.map(field => 
      new RegExp(`\\b${field}\\b`, 'i')
    );
  }

  redact(obj: any): any {
    if (obj === null || obj === undefined) {
      return obj;
    }

    if (typeof obj === 'string') {
      return this.redactString(obj);
    }

    if (Array.isArray(obj)) {
      return obj.map(item => this.redact(item));
    }

    if (typeof obj === 'object') {
      const redacted: any = {};
      for (const [key, value] of Object.entries(obj)) {
        const redactedKey = this.shouldRedactKey(key) ? '[REDACTED]' : key;
        redacted[redactedKey] = this.redact(value);
      }
      return redacted;
    }

    return obj;
  }

  private shouldRedactKey(key: string): boolean {
    const lowerKey = key.toLowerCase();
    
    // Check default blocklist
    for (const pattern of this.patterns) {
      if (pattern.test(lowerKey)) {
        return true;
      }
    }

    // Check custom rules
    for (const rule of this.rules) {
      if (rule.field === '*' || rule.field === key) {
        try {
          const regex = new RegExp(rule.pattern, 'i');
          if (regex.test(key) || regex.test(lowerKey)) {
            return true;
          }
        } catch (e) {
          console.warn('Invalid regex pattern:', rule.pattern);
        }
      }
    }

    return false;
  }

  private redactString(str: string): string {
    // Check if string contains PII patterns
    for (const rule of this.rules) {
      if (!rule.field || rule.field === '*') {
        try {
          const regex = new RegExp(rule.pattern, 'gi');
          if (regex.test(str)) {
            return str.replace(regex, rule.replacement || '[REDACTED]');
          }
        } catch (e) {
          console.warn('Invalid regex pattern:', rule.pattern);
        }
      }
    }
    return str;
  }
}

export function redact(data: any, rules: RedactionRule[] = []): any {
  const redactor = new Redactor(rules);
  return redactor.redact(data);
}
