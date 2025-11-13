class SecurityPlatformAgent < Formula
  desc "Security Platform Agent for observability and security"
  homepage "https://github.com/security-platform/agents"
  url "https://github.com/security-platform/agents/releases/download/v1.0.0/security-platform-agent-1.0.0-darwin-amd64.tar.gz"
  sha256 "placeholder-sha256"
  version "1.0.0"

  if Hardware::CPU.arm?
    url "https://github.com/security-platform/agents/releases/download/v1.0.0/security-platform-agent-1.0.0-darwin-arm64.tar.gz"
    sha256 "placeholder-sha256-arm64"
  end

  def install
    bin.install "security-platform-agent"
  end

  def post_install
    # Create config directory
    (etc/"security-platform").mkpath
    
    # Create default config if it doesn't exist
    config_file = etc/"security-platform"/"agent-config.yaml"
    unless config_file.exist?
      config_file.write <<~EOS
        control_plane_url: "https://api.securityplatform.com"
        service_name: "security-platform-agent"
        environment: "production"
      EOS
    end
  end

  plist_options :startup => true, :manual => "security-platform-agent"
  
  def plist
    <<~EOS
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
      <dict>
        <key>Label</key>
        <string>#{plist_name}</string>
        <key>ProgramArguments</key>
        <array>
          <string>#{opt_bin}/security-platform-agent</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
        <key>KeepAlive</key>
        <true/>
        <key>StandardOutPath</key>
        <string>#{var}/log/security-platform-agent.log</string>
        <key>StandardErrorPath</key>
        <string>#{var}/log/security-platform-agent.error.log</string>
        <key>EnvironmentVariables</key>
        <dict>
          <key>SECURITY_PLATFORM_CONTROL_PLANE_URL</key>
          <string>https://api.securityplatform.com</string>
        </dict>
      </dict>
      </plist>
    EOS
  end

  test do
    system "#{bin}/security-platform-agent", "--version"
  end
end
