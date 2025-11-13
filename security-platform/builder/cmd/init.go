package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initLang string
var initProject string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Scaffold agent template for chosen runtime",
	Long: `Initialize a new agent project with templates for the specified language.
Supported languages: node, python, java, dotnet, go, frontend, ios, android, desktop`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		if initLang == "" {
			return fmt.Errorf("--lang is required")
		}

		projectPath := filepath.Join(".", projectName)
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}

		fmt.Printf("Initializing %s agent project: %s\n", initLang, projectName)
		
		// Create template based on language
		switch initLang {
		case "node":
			return scaffoldNodeAgent(projectPath)
		case "python":
			return scaffoldPythonAgent(projectPath)
		case "java":
			return scaffoldJavaAgent(projectPath)
		case "dotnet":
			return scaffoldDotNetAgent(projectPath)
		case "go":
			return scaffoldGoAgent(projectPath)
		case "frontend":
			return scaffoldFrontendSDK(projectPath)
		case "ios":
			return scaffoldIOSSDK(projectPath)
		case "android":
			return scaffoldAndroidSDK(projectPath)
		case "desktop":
			return scaffoldDesktopAgent(projectPath)
		default:
			return fmt.Errorf("unsupported language: %s", initLang)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initLang, "lang", "", "language/runtime (node, python, java, dotnet, go, frontend, ios, android, desktop)")
}

func scaffoldNodeAgent(path string) error {
	// Create basic Node.js agent structure
	files := map[string]string{
		"package.json": `{
  "name": "@security-platform/node-agent",
  "version": "1.0.0",
  "description": "Security Platform Node.js Agent",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "test": "jest"
  },
  "dependencies": {
    "@opentelemetry/api": "^1.7.0",
    "@opentelemetry/sdk-trace-node": "^1.20.0",
    "@opentelemetry/exporter-otlp-http": "^0.45.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0"
  }
}`,
		"tsconfig.json": `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "declaration": true,
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true
  },
  "include": ["src/**/*"]
}`,
		"src/index.ts": `import { Agent } from './agent';

export { Agent };
export * from './types';
`,
		"src/agent.ts": `import { NodeSDK } from '@opentelemetry/sdk-trace-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { AgentConfig } from './types';

export class Agent {
  private sdk: NodeSDK | null = null;
  private config: AgentConfig;

  constructor(config: AgentConfig) {
    this.config = config;
  }

  start(): void {
    const exporter = new OTLPTraceExporter({
      url: this.config.otlpEndpoint || 'http://localhost:4318/v1/traces',
    });

    this.sdk = new NodeSDK({
      resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: this.config.serviceName || 'unknown-service',
      }),
      traceExporter: exporter,
    });

    this.sdk.start();
  }

  stop(): void {
    this.sdk?.shutdown();
  }
}
`,
		"src/types.ts": `export interface AgentConfig {
  controlPlaneUrl?: string;
  authKey?: string;
  serviceName?: string;
  otlpEndpoint?: string;
  environment?: string;
  redactionRules?: string[];
  localPolicy?: 'observe' | 'block';
  telemetryBatchSize?: number;
}
`,
	}

	return createFiles(path, files)
}

func scaffoldPythonAgent(path string) error {
	files := map[string]string{
		"setup.py": `from setuptools import setup, find_packages

setup(
    name="security-platform-agent",
    version="1.0.0",
    packages=find_packages(),
    install_requires=[
        "opentelemetry-api>=1.20.0",
        "opentelemetry-sdk>=1.20.0",
        "opentelemetry-exporter-otlp>=1.20.0",
        "fastapi>=0.100.0",
    ],
)
`,
		"pyproject.toml": `[build-system]
requires = ["setuptools>=61.0"]
build-backend = "setuptools.build_meta"

[project]
name = "security-platform-agent"
version = "1.0.0"
requires-python = ">=3.8"
`,
		"src/agent/__init__.py": `from .agent import Agent
from .types import AgentConfig

__all__ = ["Agent", "AgentConfig"]
`,
		"src/agent/agent.py": `from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from .types import AgentConfig

class Agent:
    def __init__(self, config: AgentConfig):
        self.config = config
        self.tracer_provider = None

    def start(self):
        resource = Resource.create({
            "service.name": self.config.get("service_name", "unknown-service"),
        })
        
        self.tracer_provider = TracerProvider(resource=resource)
        
        otlp_exporter = OTLPSpanExporter(
            endpoint=self.config.get("otlp_endpoint", "http://localhost:4318/v1/traces")
        )
        
        span_processor = BatchSpanProcessor(otlp_exporter)
        self.tracer_provider.add_span_processor(span_processor)
        
        trace.set_tracer_provider(self.tracer_provider)

    def stop(self):
        if self.tracer_provider:
            self.tracer_provider.shutdown()
`,
		"src/agent/types.py": `from typing import Optional, List, Literal

AgentConfig = dict[str, any]
`,
	}

	return createFiles(path, files)
}

func scaffoldJavaAgent(path string) error {
	files := map[string]string{
		"pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    
    <groupId>com.securityplatform</groupId>
    <artifactId>agent</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>
    
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    </properties>
    
    <dependencies>
        <dependency>
            <groupId>io.opentelemetry</groupId>
            <artifactId>opentelemetry-api</artifactId>
            <version>1.32.0</version>
        </dependency>
        <dependency>
            <groupId>io.opentelemetry</groupId>
            <artifactId>opentelemetry-sdk</artifactId>
            <version>1.32.0</version>
        </dependency>
        <dependency>
            <groupId>io.opentelemetry</groupId>
            <artifactId>opentelemetry-exporter-otlp</artifactId>
            <version>1.32.0</version>
        </dependency>
    </dependencies>
</project>
`,
		"src/main/java/com/securityplatform/agent/Agent.java": `package com.securityplatform.agent;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor;
import io.opentelemetry.exporter.otlp.trace.OtlpSpanExporter;

public class Agent {
    private OpenTelemetry openTelemetry;
    private AgentConfig config;

    public Agent(AgentConfig config) {
        this.config = config;
    }

    public void start() {
        OtlpSpanExporter spanExporter = OtlpSpanExporter.builder()
            .setEndpoint(config.getOtlpEndpoint())
            .build();

        SdkTracerProvider tracerProvider = SdkTracerProvider.builder()
            .addSpanProcessor(BatchSpanProcessor.builder(spanExporter).build())
            .setResource(Resource.getDefault()
                .merge(Resource.builder()
                    .put(ResourceAttributes.SERVICE_NAME, config.getServiceName())
                    .build()))
            .build();

        this.openTelemetry = OpenTelemetrySdk.builder()
            .setTracerProvider(tracerProvider)
            .build();
    }

    public void stop() {
        if (openTelemetry instanceof OpenTelemetrySdk) {
            ((OpenTelemetrySdk) openTelemetry).getSdkTracerProvider().shutdown();
        }
    }
}
`,
	}

	return createFiles(path, files)
}

func scaffoldDotNetAgent(path string) error {
	files := map[string]string{
		"SecurityPlatform.Agent.csproj": `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <LangVersion>latest</LangVersion>
  </PropertyGroup>
  
  <ItemGroup>
    <PackageReference Include="OpenTelemetry" Version="1.7.0" />
    <PackageReference Include="OpenTelemetry.Exporter.Otlp" Version="1.7.0" />
    <PackageReference Include="OpenTelemetry.Extensions.Hosting" Version="1.7.0" />
  </ItemGroup>
</Project>
`,
		"Agent.cs": `using OpenTelemetry;
using OpenTelemetry.Trace;
using OpenTelemetry.Resources;

namespace SecurityPlatform.Agent
{
    public class Agent
    {
        private TracerProvider? tracerProvider;
        private AgentConfig config;

        public Agent(AgentConfig config)
        {
            this.config = config;
        }

        public void Start()
        {
            var resourceBuilder = ResourceBuilder.CreateDefault()
                .AddService(config.ServiceName ?? "unknown-service");

            this.tracerProvider = Sdk.CreateTracerProviderBuilder()
                .SetResourceBuilder(resourceBuilder)
                .AddOtlpExporter(options =>
                {
                    options.Endpoint = new Uri(config.OtlpEndpoint ?? "http://localhost:4318/v1/traces");
                })
                .Build();
        }

        public void Stop()
        {
            tracerProvider?.Dispose();
        }
    }
}
`,
	}

	return createFiles(path, files)
}

func scaffoldGoAgent(path string) error {
	files := map[string]string{
		"go.mod": `module github.com/security-platform/go-agent

go 1.21

require (
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
)
`,
		"agent.go": `package agent

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Agent struct {
	tracerProvider *trace.TracerProvider
	config         *Config
}

type Config struct {
	ServiceName  string
	OtlpEndpoint string
	Environment  string
}

func NewAgent(config *Config) *Agent {
	return &Agent{config: config}
}

func (a *Agent) Start(ctx context.Context) error {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(a.config.OtlpEndpoint),
	)
	if err != nil {
		return err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(a.config.ServiceName),
		),
	)
	if err != nil {
		return err
	}

	a.tracerProvider = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(a.tracerProvider)
	return nil
}

func (a *Agent) Stop(ctx context.Context) error {
	if a.tracerProvider != nil {
		return a.tracerProvider.Shutdown(ctx)
	}
	return nil
}
`,
	}

	return createFiles(path, files)
}

func scaffoldFrontendSDK(path string) error {
	files := map[string]string{
		"package.json": `{
  "name": "@security-platform/frontend-sdk",
  "version": "1.0.0",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc && webpack",
    "test": "jest"
  },
  "dependencies": {
    "@opentelemetry/api": "^1.7.0",
    "@opentelemetry/sdk-trace-web": "^1.20.0",
    "@opentelemetry/exporter-otlp-http": "^0.45.0",
    "rrweb": "^2.0.0"
  }
}`,
		"src/index.ts": `import { init } from './sdk';

export { init };
export * from './types';
`,
		"src/sdk.ts": `import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { record } from 'rrweb';
import { SDKConfig } from './types';

export function init(config: SDKConfig): void {
  // Initialize OpenTelemetry
  const tracerProvider = new WebTracerProvider({
    resource: {
      serviceName: config.serviceName || 'web-app',
    },
  });

  const exporter = new OTLPTraceExporter({
    url: config.otlpEndpoint || 'http://localhost:4318/v1/traces',
  });

  tracerProvider.addSpanProcessor(new BatchSpanProcessor(exporter));
  tracerProvider.register();

  // Initialize session recording
  if (config.enableSessionRecording) {
    record({
      emit(event) {
        // Send to control plane
        fetch(config.controlPlaneUrl + '/api/sessions', {
          method: 'POST',
          body: JSON.stringify(event),
        });
      },
    });
  }
}
`,
	}

	return createFiles(path, files)
}

func scaffoldIOSSDK(path string) error {
	files := map[string]string{
		"SecurityPlatform.podspec": `Pod::Spec.new do |s|
  s.name             = 'SecurityPlatform'
  s.version          = '1.0.0'
  s.summary          = 'Security Platform iOS SDK'
  s.homepage         = 'https://github.com/security-platform/ios-sdk'
  s.license          = { :type => 'MIT' }
  s.author           = { 'Security Platform' => 'support@securityplatform.com' }
  s.source           = { :git => 'https://github.com/security-platform/ios-sdk.git', :tag => s.version }
  s.ios.deployment_target = '13.0'
  s.source_files = 'Sources/**/*'
  s.dependency 'OpenTelemetrySwift', '~> 1.0'
end
`,
		"Sources/SecurityPlatform/Agent.swift": `import Foundation
import OpenTelemetrySwift

public class Agent {
    private var tracerProvider: TracerProvider?
    private let config: AgentConfig

    public init(config: AgentConfig) {
        self.config = config
    }

    public func start() {
        // Initialize OpenTelemetry
        // Configure session recording
    }

    public func stop() {
        tracerProvider?.shutdown()
    }
}
`,
	}

	return createFiles(path, files)
}

func scaffoldAndroidSDK(path string) error {
	files := map[string]string{
		"build.gradle.kts": `plugins {
    id("com.android.library")
    id("org.jetbrains.kotlin.android")
}

android {
    namespace = "com.securityplatform.agent"
    compileSdk = 34

    defaultConfig {
        minSdk = 21
    }
}

dependencies {
    implementation("io.opentelemetry:opentelemetry-api:1.32.0")
    implementation("io.opentelemetry:opentelemetry-sdk:1.32.0")
    implementation("io.opentelemetry:opentelemetry-exporter-otlp:1.32.0")
}
`,
		"src/main/java/com/securityplatform/agent/Agent.kt": `package com.securityplatform.agent

import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.sdk.OpenTelemetrySdk

class Agent(private val config: AgentConfig) {
    private var openTelemetry: OpenTelemetry? = null

    fun start() {
        // Initialize OpenTelemetry
        // Configure session recording
    }

    fun stop() {
        openTelemetry?.let { /* shutdown */ }
    }
}
`,
	}

	return createFiles(path, files)
}

func scaffoldDesktopAgent(path string) error {
	files := map[string]string{
		"README.md": `# Desktop Agent

Cross-platform desktop agent for Windows, macOS, and Linux.
`,
		"go.mod": `module github.com/security-platform/desktop-agent

go 1.21
`,
		"main.go": `package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Start agent
	agent := NewAgent(LoadConfig())
	if err := agent.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	agent.Stop(ctx)
}
`,
	}

	return createFiles(path, files)
}

func createFiles(basePath string, files map[string]string) error {
	for relPath, content := range files {
		fullPath := filepath.Join(basePath, relPath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}
	fmt.Printf("Created project files in %s\n", basePath)
	return nil
}
