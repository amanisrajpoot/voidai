Pod::Spec.new do |spec|
  spec.name         = "SecurityPlatform"
  spec.version      = "1.0.0"
  spec.summary      = "Security Platform SDK for iOS"
  spec.description  = "Simple, automatic security observability for iOS applications"
  spec.homepage     = "https://github.com/security-platform/agent-ios"
  spec.license      = { :type => "MIT", :file => "LICENSE" }
  spec.author       = { "Security Platform" => "support@securityplatform.com" }
  spec.platform     = :ios, "13.0"
  spec.source       = { :git => "https://github.com/security-platform/agent-ios.git", :tag => "#{spec.version}" }
  spec.source_files = "SecurityPlatform/**/*.swift"
  spec.swift_version = "5.0"
  
  spec.dependency "OpenTelemetryApi", "~> 1.0"
  spec.dependency "OpenTelemetrySdk", "~> 1.0"
  spec.dependency "OpenTelemetryProtocolExporterHttp", "~> 1.0"
end
