#!/bin/bash
set -e

VERSION=${1:-"1.0.0"}
OUTPUT_DIR="./airgap-bundle-${VERSION}"

echo "Generating airgap bundle for version ${VERSION}"

# Create bundle directory
mkdir -p "${OUTPUT_DIR}"/{artifacts,images,scripts,manifests}

# Copy all built artifacts
echo "Collecting artifacts..."
cp -r dist/* "${OUTPUT_DIR}/artifacts/" 2>/dev/null || true

# Export Docker images
echo "Exporting Docker images..."
docker save observability-platform/control-plane:${VERSION} -o "${OUTPUT_DIR}/images/control-plane.tar" || true
docker save observability-platform/otel-collector:latest -o "${OUTPUT_DIR}/images/otel-collector.tar" || true
docker save observability-platform/waf:latest -o "${OUTPUT_DIR}/images/waf.tar" || true

# Create installation script
cat > "${OUTPUT_DIR}/scripts/install.sh" << 'EOF'
#!/bin/bash
set -e

echo "Installing Observability Platform (Airgap Mode)"

# Load Docker images
echo "Loading Docker images..."
docker load -i images/*.tar

# Install agents
if [ -f artifacts/*.deb ]; then
    sudo dpkg -i artifacts/*.deb
elif [ -f artifacts/*.rpm ]; then
    sudo rpm -i artifacts/*.rpm
fi

# Start services
echo "Starting services..."
docker-compose -f manifests/docker-compose.yml up -d

echo "Installation complete!"
EOF

chmod +x "${OUTPUT_DIR}/scripts/install.sh"

# Copy manifests
cp control-plane/docker-compose.yml "${OUTPUT_DIR}/manifests/"
cp control-plane/otel-collector-config.yaml "${OUTPUT_DIR}/manifests/"

# Create manifest file
cat > "${OUTPUT_DIR}/manifest.json" << EOF
{
  "version": "${VERSION}",
  "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "artifacts": $(find "${OUTPUT_DIR}/artifacts" -type f | wc -l),
  "images": $(find "${OUTPUT_DIR}/images" -type f | wc -l),
  "checksums": {}
}
EOF

# Generate checksums
echo "Generating checksums..."
find "${OUTPUT_DIR}" -type f -exec sha256sum {} \; > "${OUTPUT_DIR}/checksums.txt"

# Create tarball
echo "Creating bundle tarball..."
tar -czf "airgap-bundle-${VERSION}.tar.gz" "${OUTPUT_DIR}"

echo "Airgap bundle created: airgap-bundle-${VERSION}.tar.gz"
