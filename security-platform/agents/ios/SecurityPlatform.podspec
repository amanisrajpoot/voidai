Pod::Spec.new do |s|
  s.name             = 'SecurityPlatform'
  s.version          = '1.0.0'
  s.summary          = 'Security Platform iOS SDK for observability and security'
  s.description      = <<-DESC
  Security Platform iOS SDK provides OpenTelemetry instrumentation,
  session recording, and security monitoring for iOS applications.
  DESC
  s.homepage         = 'https://github.com/security-platform/ios-sdk'
  s.license          = { :type => 'MIT', :file => 'LICENSE' }
  s.author           = { 'Security Platform' => 'support@securityplatform.com' }
  s.source           = { :git => 'https://github.com/security-platform/ios-sdk.git', :tag => s.version.to_s }
  s.ios.deployment_target = '13.0'
  s.swift_version = '5.0'
  s.source_files = 'Sources/SecurityPlatform/**/*.swift'
  s.dependency 'OpenTelemetrySwift', '~> 1.0'
  s.dependency 'SwiftProtobuf', '~> 1.0'
end
