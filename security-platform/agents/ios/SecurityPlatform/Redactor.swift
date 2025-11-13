import Foundation

class Redactor {
    private let rules: [RedactionRule]
    private static let defaultReplacement = "***REDACTED***"

    init(rules: [RedactionRule]) {
        self.rules = rules
    }

    func redact(_ data: Any) -> Any {
        if let dict = data as? [String: Any] {
            return redactDictionary(dict)
        } else if let array = data as? [Any] {
            return array.map { redact($0) }
        } else if let str = data as? String {
            return redactString(str)
        }
        return data
    }

    private func redactDictionary(_ dict: [String: Any]) -> [String: Any] {
        var result: [String: Any] = [:]

        for (key, value) in dict {
            var shouldRedact = false

            for rule in rules {
                if let field = rule.field, field.lowercased() == key.lowercased() {
                    shouldRedact = true
                    break
                }
                if let pattern = rule.pattern {
                    let regex = try? NSRegularExpression(pattern: pattern, options: .caseInsensitive)
                    if let range = NSRange(key, in: key) as NSRange?,
                       regex?.firstMatch(in: key, range: range) != nil {
                        shouldRedact = true
                        break
                    }
                }
            }

            if shouldRedact {
                let replacement = rules.first(where: { $0.field == key || $0.pattern != nil })?.replacement ?? Self.defaultReplacement
                result[key] = replacement
            } else {
                result[key] = redact(value)
            }
        }

        return result
    }

    private func redactString(_ str: String) -> String {
        var result = str
        for rule in rules {
            if let pattern = rule.pattern {
                let replacement = rule.replacement ?? Self.defaultReplacement
                if let regex = try? NSRegularExpression(pattern: pattern, options: .caseInsensitive) {
                    result = regex.stringByReplacingMatches(in: result,
                                                             options: [],
                                                             range: NSRange(location: 0, length: result.utf16.count),
                                                             withTemplate: replacement)
                }
            }
        }
        return result
    }
}
