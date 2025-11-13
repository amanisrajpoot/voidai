// swift-tools-version: 5.7
import PackageDescription

let package = Package(
    name: "SecurityPlatform",
    platforms: [
        .iOS(.v13)
    ],
    products: [
        .library(
            name: "SecurityPlatform",
            targets: ["SecurityPlatform"]
        )
    ],
    dependencies: [
        .package(url: "https://github.com/open-telemetry/opentelemetry-swift", from: "1.0.0")
    ],
    targets: [
        .target(
            name: "SecurityPlatform",
            dependencies: [
                .product(name: "OpenTelemetryApi", package: "opentelemetry-swift"),
                .product(name: "OpenTelemetrySdk", package: "opentelemetry-swift")
            ]
        )
    ]
)
