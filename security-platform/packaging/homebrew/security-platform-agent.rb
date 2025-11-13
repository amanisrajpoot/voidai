class SecurityPlatformAgent < Formula
  desc "Security observability agent for desktop systems"
  homepage "https://github.com/securityplatform/agent"
  url "https://github.com/securityplatform/agent/releases/download/v1.0.0/security-platform-agent-darwin-amd64.tar.gz"
  sha256 "abc123def456..." # Update with actual checksum
  version "1.0.0"

  depends_on "go" => :build

  def install
    bin.install "security-platform-agent"
    etc.install "config.yaml.example" => "security-platform/config.yaml.example"
    
    # Install LaunchDaemon plist
    (prefix/"Library/LaunchDaemons").install "com.securityplatform.agent.plist"
  end

  def post_install
    # Create config directory
    (etc/"security-platform").mkpath
    
    # Copy example config if it doesn't exist
    unless (etc/"security-platform/config.yaml").exist?
      cp etc/"security-platform/config.yaml.example", etc/"security-platform/config.yaml"
    end
  end

  plist_options :startup => true
  plist_options :manual => "security-platform-agent -config=#{etc}/security-platform/config.yaml"

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
          <string>#{bin}/security-platform-agent</string>
          <string>-config</string>
          <string>#{etc}/security-platform/config.yaml</string>
          <string>-service</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
        <key>KeepAlive</key>
        <true/>
        <key>StandardOutPath</key>
        <string>#{var}/log/security-platform-agent.log</string>
        <key>StandardErrorPath</key>
        <string>#{var}/log/security-platform-agent.error.log</string>
      </dict>
      </plist>
    EOS
  end

  test do
    system "#{bin}/security-platform-agent", "--version"
  end
end
