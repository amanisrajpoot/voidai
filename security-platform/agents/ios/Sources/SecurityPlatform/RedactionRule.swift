import Foundation

/// Redaction rule configuration.
public struct RedactionRule {
    public var pattern: String
    public var replacement: String?
    public var field: String?

    public init(pattern: String, replacement: String? = nil, field: String? = nil) {
        self.pattern = pattern
        self.replacement = replacement
        self.field = field
    }
}
