Pod::Spec.new do |spec|
  spec.name         = "ObservabilityAgent"
  spec.version      = "1.0.0"
  spec.summary      = "Security Observability SDK for iOS"
  spec.description  = <<-DESC
    Security observability SDK for iOS applications. Provides telemetry collection,
    session replay, and security event tracking with automatic PII redaction.
  DESC
  
  spec.homepage     = "https://github.com/observability-platform/agent-ios"
  spec.license      = { :type => "Apache-2.0", :file => "LICENSE" }
  spec.author       = { "Observability Platform" => "support@observability.platform" }
  
  spec.platform     = :ios, "13.0"
  spec.swift_version = "5.0"
  
  spec.source       = { :git => "https://github.com/observability-platform/agent-ios.git", :tag => "#{spec.version}" }
  spec.source_files = "Sources/**/*.swift"
  
  spec.dependency "OpenTelemetrySwift", "~> 1.0"
  
  spec.frameworks = "Foundation", "UIKit"
end
