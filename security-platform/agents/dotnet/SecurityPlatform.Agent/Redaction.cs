using System.Text.RegularExpressions;

namespace SecurityPlatform.Agent;

/// <summary>
/// PII Redaction utility
/// </summary>
public class Redaction
{
    private static readonly List<Regex> DefaultPatterns = new()
    {
        new Regex(@"(?i)(password|passwd|pwd)\s*[:=]\s*([^\s,}]+)", RegexOptions.IgnoreCase),
        new Regex(@"(?i)(api[_-]?key|apikey)\s*[:=]\s*([^\s,}]+)", RegexOptions.IgnoreCase),
        new Regex(@"(?i)(token|bearer)\s*[:=]\s*([^\s,}]+)", RegexOptions.IgnoreCase),
        new Regex(@"\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b"), // Credit card
        new Regex(@"\b\d{3}-\d{2}-\d{4}\b"), // SSN
    };

    private readonly List<Regex> _patterns;

    public Redaction(List<string>? customRules = null)
    {
        _patterns = new List<Regex>(DefaultPatterns);
        if (customRules != null)
        {
            foreach (var rule in customRules)
            {
                try
                {
                    _patterns.Add(new Regex(rule, RegexOptions.IgnoreCase));
                }
                catch
                {
                    // Invalid regex, skip
                }
            }
        }
    }

    public string Redact(string? input)
    {
        if (string.IsNullOrEmpty(input))
        {
            return input ?? string.Empty;
        }

        var result = input;
        foreach (var pattern in _patterns)
        {
            result = pattern.Replace(result, "[REDACTED]");
        }
        return result;
    }

    public Dictionary<string, object> Redact(Dictionary<string, object>? data)
    {
        if (data == null)
        {
            return new Dictionary<string, object>();
        }

        var redacted = new Dictionary<string, object>();
        foreach (var entry in data)
        {
            if (ShouldRedactKey(entry.Key))
            {
                redacted[entry.Key] = "[REDACTED]";
            }
            else if (entry.Value is string strValue)
            {
                redacted[entry.Key] = Redact(strValue);
            }
            else if (entry.Value is Dictionary<string, object> dictValue)
            {
                redacted[entry.Key] = Redact(dictValue);
            }
            else
            {
                redacted[entry.Key] = entry.Value;
            }
        }
        return redacted;
    }

    private static bool ShouldRedactKey(string key)
    {
        var lowerKey = key.ToLowerInvariant();
        return lowerKey.Contains("password", StringComparison.OrdinalIgnoreCase) ||
               lowerKey.Contains("passwd", StringComparison.OrdinalIgnoreCase) ||
               lowerKey.Contains("pwd", StringComparison.OrdinalIgnoreCase) ||
               lowerKey.Contains("api_key", StringComparison.OrdinalIgnoreCase) ||
               lowerKey.Contains("apikey", StringComparison.OrdinalIgnoreCase) ||
               lowerKey.Contains("token", StringComparison.OrdinalIgnoreCase) ||
               lowerKey.Contains("secret", StringComparison.OrdinalIgnoreCase);
    }
}
