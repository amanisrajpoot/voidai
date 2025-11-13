using System;
using System.Collections;
using System.Collections.Generic;
using System.Text.RegularExpressions;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Redacts PII from data structures.
    /// </summary>
    public class Redactor
    {
        private static readonly List<string> DefaultBlocklist = new List<string>
        {
            "password", "passwd", "pwd", "secret", "token", "api_key", "apikey",
            "credit_card", "card_number", "cvv", "ssn", "social_security", "pin", "passcode"
        };

        private readonly List<Regex> _patterns = new List<Regex>();

        public Redactor(List<RedactionRule> rules)
        {
            // Add default patterns
            foreach (var field in DefaultBlocklist)
            {
                _patterns.Add(new Regex($@".*{Regex.Escape(field)}.*", RegexOptions.IgnoreCase));
            }

            // Add custom rules
            if (rules != null)
            {
                foreach (var rule in rules)
                {
                    if (!string.IsNullOrEmpty(rule.Pattern))
                    {
                        _patterns.Add(new Regex(rule.Pattern, RegexOptions.IgnoreCase));
                    }
                }
            }
        }

        public object? Redact(object? data)
        {
            if (data == null)
            {
                return null;
            }

            if (data is IDictionary<string, object> dict)
            {
                return RedactDictionary(dict);
            }
            else if (data is IList list)
            {
                return RedactList(list);
            }
            else if (data is string str)
            {
                return RedactString(str);
            }

            return data;
        }

        private Dictionary<string, object?> RedactDictionary(IDictionary<string, object> dict)
        {
            var redacted = new Dictionary<string, object?>();
            foreach (var entry in dict)
            {
                if (ShouldRedact(entry.Key))
                {
                    redacted[entry.Key] = "***REDACTED***";
                }
                else
                {
                    redacted[entry.Key] = Redact(entry.Value);
                }
            }
            return redacted;
        }

        private List<object?> RedactList(IList list)
        {
            var redacted = new List<object?>();
            foreach (var item in list)
            {
                redacted.Add(Redact(item));
            }
            return redacted;
        }

        private string RedactString(string str)
        {
            // Simple string redaction - could be enhanced with regex matching
            return str;
        }

        private bool ShouldRedact(string fieldName)
        {
            foreach (var pattern in _patterns)
            {
                if (pattern.IsMatch(fieldName))
                {
                    return true;
                }
            }
            return false;
        }
    }
}
