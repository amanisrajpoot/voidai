class SecurityPlatformAgent < Formula
  desc "Security Platform Desktop Agent"
  homepage "https://github.com/security-platform/desktop-agent"
  url "https://github.com/security-platform/desktop-agent/releases/download/v1.0.0/security-platform-agent-darwin-amd64.tar.gz"
  sha256 "placeholder-sha256"
  version "1.0.0"

  def install
    bin.install "security-platform-agent"
    etc.install "config/agent-config.yaml" => "security-platform/agent-config.yaml"
  end

  def post_install
    # Create config directory
    (etc/"security-platform").mkpath
    
    # Create default config if it doesn't exist
    config_file = etc/"security-platform/agent-config.yaml"
    unless config_file.exist?
      config_file.write <<~EOS
        control_plane_url: "https://api.securityplatform.com"
        auth_key: ""
        service_name: "desktop-agent"
        environment: "production"
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
