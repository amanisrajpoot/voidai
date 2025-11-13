import Foundation

/// PII Redaction utility
public class Redaction {
    private let patterns: [NSRegularExpression]
    
    private static let defaultPatterns: [String] = [
        "(?i)(password|passwd|pwd)\\s*[:=]\\s*([^\\s,}]+)",
        "(?i)(api[_-]?key|apikey)\\s*[:=]\\s*([^\\s,}]+)",
        "(?i)(token|bearer)\\s*[:=]\\s*([^\\s,}]+)",
        "\\b\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}\\b", // Credit card
        "\\b\\d{3}-\\d{2}-\\d{4}\\b" // SSN
    ]
    
    public init(customRules: [String]? = nil) {
        var compiledPatterns: [NSRegularExpression] = []
        
        for pattern in Self.defaultPatterns {
            if let regex = try? NSRegularExpression(pattern: pattern, options: .caseInsensitive) {
                compiledPatterns.append(regex)
            }
        }
        
        if let customRules = customRules {
            for rule in customRules {
                if let regex = try? NSRegularExpression(pattern: rule, options: .caseInsensitive) {
                    compiledPatterns.append(regex)
                }
            }
        }
        
        self.patterns = compiledPatterns
    }
    
    public func redact(_ input: String) -> String {
        guard !input.isEmpty else { return input }
        
        var result = input
        for pattern in patterns {
            result = pattern.stringByReplacingMatches(
                in: result,
                options: [],
                range: NSRange(location: 0, length: result.utf16.count),
                withTemplate: "[REDACTED]"
            )
        }
        return result
    }
    
    public func redact(_ data: [String: Any]) -> [String: Any] {
        var redacted: [String: Any] = [:]
        
        for (key, value) in data {
            if shouldRedactKey(key) {
                redacted[key] = "[REDACTED]"
            } else if let stringValue = value as? String {
                redacted[key] = redact(stringValue)
            } else if let dictValue = value as? [String: Any] {
                redacted[key] = redact(dictValue)
            } else {
                redacted[key] = value
            }
        }
        
        return redacted
    }
    
    private func shouldRedactKey(_ key: String) -> Bool {
        let lowerKey = key.lowercased()
        return lowerKey.contains("password") ||
               lowerKey.contains("passwd") ||
               lowerKey.contains("pwd") ||
               lowerKey.contains("api_key") ||
               lowerKey.contains("apikey") ||
               lowerKey.contains("token") ||
               lowerKey.contains("secret")
    }
}
