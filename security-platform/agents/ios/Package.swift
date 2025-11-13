// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "SecurityPlatform",
    platforms: [
        .iOS(.v13),
        .macOS(.v10_15),
        .watchOS(.v6),
        .tvOS(.v13)
    ],
    products: [
        .library(
            name: "SecurityPlatform",
            targets: ["SecurityPlatform"]
        ),
    ],
    dependencies: [
        .package(url: "https://github.com/open-telemetry/opentelemetry-swift", from: "1.9.0"),
    ],
    targets: [
        .target(
            name: "SecurityPlatform",
            dependencies: [
                .product(name: "OpenTelemetryApi", package: "opentelemetry-swift"),
                .product(name: "OpenTelemetrySdk", package: "opentelemetry-swift"),
                .product(name: "ResourceExtension", package: "opentelemetry-swift"),
                .product(name: "StdoutExporter", package: "opentelemetry-swift"),
            ]
        ),
        .testTarget(
            name: "SecurityPlatformTests",
            dependencies: ["SecurityPlatform"]
        ),
    ]
)
