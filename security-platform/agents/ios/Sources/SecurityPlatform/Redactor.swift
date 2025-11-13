import Foundation

/// Redacts PII from data structures.
public class Redactor {
    private static let defaultBlocklist = [
        "password", "passwd", "pwd", "secret", "token", "api_key", "apikey",
        "credit_card", "card_number", "cvv", "ssn", "social_security", "pin", "passcode"
    ]

    private var patterns: [String]

    public init(rules: [RedactionRule]) {
        self.patterns = Self.defaultBlocklist
        // Add custom patterns from rules
        for rule in rules {
            if !rule.pattern.isEmpty {
                patterns.append(rule.pattern)
            }
        }
    }

    public func redact(_ data: Any) -> Any {
        if let dict = data as? [String: Any] {
            return redactDictionary(dict)
        } else if let array = data as? [Any] {
            return redactArray(array)
        } else if let str = data as? String {
            return redactString(str)
        }
        return data
    }

    private func redactDictionary(_ dict: [String: Any]) -> [String: Any] {
        var redacted: [String: Any] = [:]
        for (key, value) in dict {
            if shouldRedact(key) {
                redacted[key] = "***REDACTED***"
            } else {
                redacted[key] = redact(value)
            }
        }
        return redacted
    }

    private func redactArray(_ array: [Any]) -> [Any] {
        return array.map { redact($0) }
    }

    private func redactString(_ str: String) -> String {
        // Simple string redaction
        return str
    }

    private func shouldRedact(_ fieldName: String) -> Bool {
        let lowercased = fieldName.lowercased()
        for pattern in patterns {
            if lowercased.contains(pattern.lowercased()) {
                return true
            }
        }
        return false
    }
}
