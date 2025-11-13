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
        private const string DefaultReplacement = "***REDACTED***";
        private readonly List<RedactionRule> _rules;

        public Redactor(List<RedactionRule> rules)
        {
            _rules = rules ?? new List<RedactionRule>();
        }

        public object? Redact(object? data)
        {
            if (data == null)
                return null;

            return data switch
            {
                Dictionary<string, object> dict => RedactDictionary(dict),
                IDictionary<string, object> dict => RedactDictionary(dict.ToDictionary(kvp => kvp.Key, kvp => kvp.Value)),
                System.Collections.IList list => RedactList(list),
                string str => RedactString(str),
                _ => data
            };
        }

        private Dictionary<string, object> RedactDictionary(Dictionary<string, object> dict)
        {
            var result = new Dictionary<string, object>();
            foreach (var kvp in dict)
            {
                var key = kvp.Key;
                var value = kvp.Value;

                bool shouldRedact = _rules.Any(rule =>
                    (!string.IsNullOrEmpty(rule.Field) && rule.Field.Equals(key, StringComparison.OrdinalIgnoreCase)) ||
                    (!string.IsNullOrEmpty(rule.Pattern) && Regex.IsMatch(key, rule.Pattern, RegexOptions.IgnoreCase)));

                result[key] = shouldRedact ? DefaultReplacement : Redact(value);
            }
            return result;
        }

        private Dictionary<string, object> RedactDictionary(IDictionary<string, object> dict)
        {
            return RedactDictionary(dict.ToDictionary(kvp => kvp.Key, kvp => kvp.Value));
        }

        private System.Collections.IList RedactList(System.Collections.IList list)
        {
            var result = new List<object?>();
            foreach (var item in list)
            {
                result.Add(Redact(item));
            }
            return result;
        }

        private string RedactString(string str)
        {
            foreach (var rule in _rules)
            {
                if (!string.IsNullOrEmpty(rule.Pattern))
                {
                    var replacement = rule.Replacement ?? DefaultReplacement;
                    str = Regex.Replace(str, rule.Pattern, replacement, RegexOptions.IgnoreCase);
                }
            }
            return str;
        }
    }
}
