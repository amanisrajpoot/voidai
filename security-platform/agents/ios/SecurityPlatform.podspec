Pod::Spec.new do |s|
  s.name             = 'SecurityPlatform'
  s.version          = '1.0.0'
  s.summary          = 'Security Platform Agent for iOS'
  s.description      = 'Security observability and protection SDK for iOS applications'
  s.homepage         = 'https://github.com/security-platform/agents'
  s.license          = { :type => 'MIT', :file => 'LICENSE' }
  s.author           = { 'Security Platform' => 'support@securityplatform.com' }
  s.source           = { :git => 'https://github.com/security-platform/agents.git', :tag => s.version.to_s }
  s.ios.deployment_target = '13.0'
  s.swift_version = '5.0'
  s.source_files = 'SecurityPlatform/**/*.swift'
  s.dependency 'OpenTelemetrySwift', '~> 1.0'
end
