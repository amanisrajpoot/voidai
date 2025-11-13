class SecurityPlatformAgent < Formula
  desc "Security Platform Desktop Agent for observability and security"
  homepage "https://github.com/security-platform/desktop-agent"
  url "https://github.com/security-platform/desktop-agent/releases/download/v1.0.0/security-platform-agent-darwin-amd64.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  version "1.0.0"

  def install
    bin.install "security-platform-agent"
    etc.install "config.yaml" => "security-platform/config.yaml" if File.exist?("config.yaml")
  end

  def post_install
    # Create config directory if it doesn't exist
    config_dir = etc/"security-platform"
    config_dir.mkpath unless config_dir.exist?
    
    # Create default config if it doesn't exist
    default_config = config_dir/"config.yaml"
    unless default_config.exist?
      default_config.write <<~EOS
        service_name: desktop-agent
        version: 1.0.0
        environment: production
        otlp_endpoint: http://localhost:4318/v1/traces
        local_policy: observe
      EOS
    end
  end

  service do
    run [opt_bin/"security-platform-agent", "-service"]
    keep_alive true
    log_path var/"log/security-platform-agent.log"
    error_log_path var/"log/security-platform-agent.error.log"
  end

  test do
    system "#{bin}/security-platform-agent", "--version"
  end
end
