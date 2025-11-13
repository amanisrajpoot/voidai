using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;

namespace SecurityPlatform.Agent
{
    /// <summary>
    /// Redacts PII from data structures.
    /// </summary>
    public class Redactor
    {
        private readonly List<string> _blocklist;
        private const string Redacted = "[REDACTED]";

        private static readonly List<string> DefaultBlocklist = new List<string>
        {
            "password", "passwd", "pwd", "secret", "token",
            "api_key", "apikey", "credit_card", "card_number",
            "cvv", "ssn", "social_security", "pin", "passcode"
        };

        public Redactor(List<string>? blocklist = null)
        {
            _blocklist = blocklist ?? DefaultBlocklist;
        }

        public object? Redact(object? data)
        {
            if (data == null)
            {
                return null;
            }

            return data switch
            {
                Dictionary<string, object?> dict => RedactDictionary(dict),
                IDictionary<string, object?> dict => RedactDictionary(dict),
                IEnumerable<object?> list => RedactList(list),
                string str => RedactString(str),
                _ => data
            };
        }

        private Dictionary<string, object?> RedactDictionary(IDictionary<string, object?> dict)
        {
            var redacted = new Dictionary<string, object?>();
            foreach (var kvp in dict)
            {
                if (ShouldRedact(kvp.Key))
                {
                    redacted[kvp.Key] = Redacted;
                }
                else
                {
                    redacted[kvp.Key] = Redact(kvp.Value);
                }
            }
            return redacted;
        }

        private List<object?> RedactList(IEnumerable<object?> list)
        {
            return list.Select(Redact).ToList();
        }

        private string RedactString(string str)
        {
            // Simple string redaction - could be enhanced with regex patterns
            return str;
        }

        private bool ShouldRedact(string? key)
        {
            if (string.IsNullOrEmpty(key))
            {
                return false;
            }

            var lowerKey = key.ToLowerInvariant();
            foreach (var pattern in _blocklist)
            {
                if (pattern.Contains("*"))
                {
                    var regex = "^" + Regex.Escape(pattern).Replace("\\*", ".*") + "$";
                    if (Regex.IsMatch(lowerKey, regex, RegexOptions.IgnoreCase))
                    {
                        return true;
                    }
                }
                else if (lowerKey.Contains(pattern.ToLowerInvariant()))
                {
                    return true;
                }
            }
            return false;
        }
    }
}
