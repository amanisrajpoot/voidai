package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Scaffolder struct {
	projectPath string
	lang        string
}

func NewScaffolder(projectPath, lang string) *Scaffolder {
	return &Scaffolder{
		projectPath: projectPath,
		lang:        lang,
	}
}

func (s *Scaffolder) Scaffold() error {
	switch s.lang {
	case "node":
		return s.scaffoldNode()
	case "python":
		return s.scaffoldPython()
	case "java":
		return s.scaffoldJava()
	case "dotnet":
		return s.scaffoldDotNet()
	case "go":
		return s.scaffoldGo()
	case "ruby":
		return s.scaffoldRuby()
	case "frontend":
		return s.scaffoldFrontend()
	case "ios":
		return s.scaffoldIOS()
	case "android":
		return s.scaffoldAndroid()
	case "desktop-windows":
		return s.scaffoldDesktopWindows()
	case "desktop-macos":
		return s.scaffoldDesktopMacOS()
	case "desktop-linux":
		return s.scaffoldDesktopLinux()
	default:
		return fmt.Errorf("unsupported language: %s", s.lang)
	}
}

func (s *Scaffolder) writeFile(path string, content string) error {
	fullPath := filepath.Join(s.projectPath, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, []byte(content), 0644)
}

func (s *Scaffolder) writeTemplate(path string, tmpl string, data interface{}) error {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	fullPath := filepath.Join(s.projectPath, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, data)
}

func (s *Scaffolder) scaffoldNode() error {
	files := map[string]string{
		"package.json": `{
  "name": "observability-agent-node",
  "version": "1.0.0",
  "description": "Security Observability Agent for Node.js",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "test": "jest"
  },
  "dependencies": {
    "@opentelemetry/api": "^1.7.0",
    "@opentelemetry/sdk-trace-node": "^1.15.0",
    "@opentelemetry/exporter-otlp-http": "^0.45.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "jest": "^29.0.0"
  }
}`,
		"tsconfig.json": `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}`,
		"src/index.ts": `import { NodeSDK } from '@opentelemetry/sdk-trace-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { ExpressInstrumentation } from '@opentelemetry/instrumentation-express';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { registerInstrumentations } from '@opentelemetry/instrumentation';

export interface AgentConfig {
  controlPlaneUrl: string;
  authKey: string;
  serviceName: string;
  environment?: string;
  otlpEndpoint?: string;
  batchSize?: number;
  redactionRules?: string[];
}

export class ObservabilityAgent {
  private sdk: NodeSDK | null = null;

  init(config: AgentConfig): void {
    const resource = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: config.serviceName,
      [SemanticResourceAttributes.SERVICE_VERSION]: '1.0.0',
      [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: config.environment || 'production',
    });

    const exporter = new OTLPTraceExporter({
      url: config.otlpEndpoint || `${config.controlPlaneUrl}/v1/traces`,
      headers: {
        'Authorization': `Bearer ${config.authKey}`,
      },
    });

    this.sdk = new NodeSDK({
      resource,
      traceExporter: exporter,
      instrumentations: [
        new HttpInstrumentation(),
        new ExpressInstrumentation(),
      ],
    });

    this.sdk.start();
  }

  shutdown(): void {
    if (this.sdk) {
      this.sdk.shutdown();
    }
  }
}

export default ObservabilityAgent;
`,
		"src/middleware.ts": `import { Request, Response, NextFunction } from 'express';
import { context, trace } from '@opentelemetry/api';
import { redact } from './redaction';

export function observabilityMiddleware(req: Request, res: Response, next: NextFunction): void {
  const tracer = trace.getTracer('observability-agent');
  const span = tracer.startSpan(`${req.method} ${req.path}`);

  span.setAttributes({
    'http.method': req.method,
    'http.url': req.url,
    'http.route': req.route?.path || req.path,
    'user.id': req.user?.id || '',
  });

  // Capture request body (redacted)
  if (req.body) {
    const redactedBody = redact(req.body);
    span.setAttribute('http.request.body', JSON.stringify(redactedBody));
  }

  res.on('finish', () => {
    span.setAttribute('http.status_code', res.statusCode);
    span.end();
  });

  context.with(trace.setSpan(context.active(), span), () => {
    next();
  });
}
`,
		"src/redaction.ts": `const DEFAULT_REDACTION_PATTERNS = [
  /password/i,
  /passwd/i,
  /pwd/i,
  /secret/i,
  /token/i,
  /api[_-]?key/i,
  /auth[_-]?token/i,
  /credit[_-]?card/i,
  /card[_-]?number/i,
  /cvv/i,
  /ssn/i,
  /social[_-]?security/i,
  /pin/i,
];

export function redact(obj: any, patterns: RegExp[] = DEFAULT_REDACTION_PATTERNS): any {
  if (obj === null || obj === undefined) {
    return obj;
  }

  if (typeof obj === 'string') {
    for (const pattern of patterns) {
      if (pattern.test(obj)) {
        return '[REDACTED]';
      }
    }
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(item => redact(item, patterns));
  }

  if (typeof obj === 'object') {
    const redacted: any = {};
    for (const [key, value] of Object.entries(obj)) {
      let shouldRedact = false;
      for (const pattern of patterns) {
        if (pattern.test(key)) {
          shouldRedact = true;
          break;
        }
      }
      redacted[key] = shouldRedact ? '[REDACTED]' : redact(value, patterns);
    }
    return redacted;
  }

  return obj;
}
`,
		"README.md": `# Node.js Observability Agent

Security observability agent for Node.js applications.

## Installation

\`\`\`bash
npm install observability-agent-node
\`\`\`

## Usage

\`\`\`typescript
import { ObservabilityAgent } from 'observability-agent-node';
import express from 'express';
import { observabilityMiddleware } from 'observability-agent-node/middleware';

const app = express();
const agent = new ObservabilityAgent();

agent.init({
  controlPlaneUrl: process.env.CONTROL_PLANE_URL || 'https://control.example.com',
  authKey: process.env.AUTH_KEY || '',
  serviceName: 'my-service',
  environment: process.env.NODE_ENV || 'production',
});

app.use(observabilityMiddleware);

app.listen(3000);
\`\`\`
`,
	}

	for path, content := range files {
		if err := s.writeFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scaffolder) scaffoldPython() error {
	files := map[string]string{
		"setup.py": `from setuptools import setup, find_packages

setup(
    name="observability-agent-python",
    version="1.0.0",
    description="Security Observability Agent for Python",
    packages=find_packages(),
    install_requires=[
        "opentelemetry-api>=1.20.0",
        "opentelemetry-sdk>=1.20.0",
        "opentelemetry-exporter-otlp-proto-http>=1.20.0",
        "opentelemetry-instrumentation-fastapi>=0.41b0",
        "opentelemetry-instrumentation-django>=0.41b0",
    ],
)
`,
		"observability_agent/__init__.py": `from .agent import ObservabilityAgent
from .middleware import observability_middleware

__all__ = ['ObservabilityAgent', 'observability_middleware']
`,
		"observability_agent/agent.py": `import os
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.semantic_conventions.resource import ResourceAttributes

class ObservabilityAgent:
    def __init__(self, config: dict):
        self.config = config
        self.tracer_provider = None
        
    def init(self):
        resource = Resource.create({
            ResourceAttributes.SERVICE_NAME: self.config.get('service_name', 'unknown'),
            ResourceAttributes.SERVICE_VERSION: self.config.get('service_version', '1.0.0'),
            ResourceAttributes.DEPLOYMENT_ENVIRONMENT: self.config.get('environment', 'production'),
        })
        
        self.tracer_provider = TracerProvider(resource=resource)
        
        otlp_endpoint = self.config.get('otlp_endpoint') or f"{self.config['control_plane_url']}/v1/traces"
        otlp_exporter = OTLPSpanExporter(
            endpoint=otlp_endpoint,
            headers={
                'Authorization': f"Bearer {self.config['auth_key']}"
            }
        )
        
        span_processor = BatchSpanProcessor(otlp_exporter)
        self.tracer_provider.add_span_processor(span_processor)
        
        trace.set_tracer_provider(self.tracer_provider)
        
    def shutdown(self):
        if self.tracer_provider:
            self.tracer_provider.shutdown()
`,
		"observability_agent/middleware.py": `from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware
from opentelemetry import trace
from .redaction import redact

tracer = trace.get_tracer(__name__)

class ObservabilityMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        span = tracer.start_span(f"{request.method} {request.url.path}")
        
        span.set_attribute("http.method", request.method)
        span.set_attribute("http.url", str(request.url))
        span.set_attribute("http.route", request.url.path)
        
        # Capture request body (redacted)
        if request.method in ["POST", "PUT", "PATCH"]:
            body = await request.body()
            if body:
                try:
                    import json
                    body_obj = json.loads(body)
                    redacted_body = redact(body_obj)
                    span.set_attribute("http.request.body", json.dumps(redacted_body))
                except:
                    pass
        
        response = await call_next(request)
        
        span.set_attribute("http.status_code", response.status_code)
        span.end()
        
        return response

def observability_middleware(app):
    return ObservabilityMiddleware(app)
`,
		"observability_agent/redaction.py": `import re
from typing import Any, Dict, List

DEFAULT_REDACTION_PATTERNS = [
    re.compile(r'password', re.IGNORECASE),
    re.compile(r'passwd', re.IGNORECASE),
    re.compile(r'pwd', re.IGNORECASE),
    re.compile(r'secret', re.IGNORECASE),
    re.compile(r'token', re.IGNORECASE),
    re.compile(r'api[_-]?key', re.IGNORECASE),
    re.compile(r'auth[_-]?token', re.IGNORECASE),
    re.compile(r'credit[_-]?card', re.IGNORECASE),
    re.compile(r'card[_-]?number', re.IGNORECASE),
    re.compile(r'cvv', re.IGNORECASE),
    re.compile(r'ssn', re.IGNORECASE),
    re.compile(r'social[_-]?security', re.IGNORECASE),
    re.compile(r'pin', re.IGNORECASE),
]

def redact(obj: Any, patterns: List[re.Pattern] = None) -> Any:
    if patterns is None:
        patterns = DEFAULT_REDACTION_PATTERNS
    
    if obj is None:
        return obj
    
    if isinstance(obj, str):
        for pattern in patterns:
            if pattern.search(obj):
                return '[REDACTED]'
        return obj
    
    if isinstance(obj, list):
        return [redact(item, patterns) for item in obj]
    
    if isinstance(obj, dict):
        redacted = {}
        for key, value in obj.items():
            should_redact = any(pattern.search(key) for pattern in patterns)
            redacted[key] = '[REDACTED]' if should_redact else redact(value, patterns)
        return redacted
    
    return obj
`,
		"README.md": `# Python Observability Agent

Security observability agent for Python applications (FastAPI, Django, Flask).

## Installation

\`\`\`bash
pip install observability-agent-python
\`\`\`

## Usage (FastAPI)

\`\`\`python
from fastapi import FastAPI
from observability_agent import ObservabilityAgent, observability_middleware

app = FastAPI()
agent = ObservabilityAgent({
    'control_plane_url': os.getenv('CONTROL_PLANE_URL', 'https://control.example.com'),
    'auth_key': os.getenv('AUTH_KEY', ''),
    'service_name': 'my-service',
    'environment': os.getenv('ENVIRONMENT', 'production'),
})

agent.init()
app.add_middleware(observability_middleware)

@app.get("/")
def read_root():
    return {"Hello": "World"}
\`\`\`
`,
	}

	for path, content := range files {
		if err := s.writeFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scaffolder) scaffoldGo() error {
	files := map[string]string{
		"go.mod": `module observability-agent-go

go 1.21

require (
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
)
`,
		"agent.go": `package main

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Config struct {
	ControlPlaneURL string
	AuthKey        string
	ServiceName    string
	Environment    string
	OTLPEndpoint   string
}

type Agent struct {
	config      Config
	tp          *trace.TracerProvider
}

func NewAgent(config Config) *Agent {
	return &Agent{config: config}
}

func (a *Agent) Init() error {
	otlpEndpoint := a.config.OTLPEndpoint
	if otlpEndpoint == "" {
		otlpEndpoint = a.config.ControlPlaneURL + "/v1/traces"
	}

	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Bearer " + a.config.AuthKey,
		}),
	)
	if err != nil {
		return err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(a.config.ServiceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String(a.config.Environment),
		),
	)
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	a.tp = tp

	return nil
}

func (a *Agent) Shutdown() error {
	if a.tp != nil {
		return a.tp.Shutdown(context.Background())
	}
	return nil
}
`,
		"middleware.go": `package main

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ObservabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("observability-agent")
		ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path)

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.route", r.URL.Path),
		)

		next.ServeHTTP(w, r.WithContext(ctx))

		span.End()
	})
}
`,
		"README.md": `# Go Observability Agent

Security observability agent for Go applications.

## Installation

\`\`\`bash
go get github.com/observability-platform/agent-go
\`\`\`

## Usage

\`\`\`go
package main

import (
	"log"
	"net/http"
	"os"
	
	agent "github.com/observability-platform/agent-go"
)

func main() {
	cfg := agent.Config{
		ControlPlaneURL: os.Getenv("CONTROL_PLANE_URL"),
		AuthKey:        os.Getenv("AUTH_KEY"),
		ServiceName:    "my-service",
		Environment:    os.Getenv("ENVIRONMENT"),
	}

	ag := agent.NewAgent(cfg)
	if err := ag.Init(); err != nil {
		log.Fatal(err)
	}
	defer ag.Shutdown()

	http.Handle("/", agent.ObservabilityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
\`\`\`
`,
	}

	for path, content := range files {
		if err := s.writeFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scaffolder) scaffoldFrontend() error {
	files := map[string]string{
		"package.json": `{
  "name": "@observability-platform/frontend-sdk",
  "version": "1.0.0",
  "description": "Security Observability SDK for Web Frontends",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc && webpack",
    "test": "jest"
  },
  "dependencies": {
    "@opentelemetry/api": "^1.7.0",
    "@opentelemetry/sdk-trace-web": "^1.15.0",
    "@opentelemetry/exporter-otlp-http": "^0.45.0",
    "rrweb": "^2.0.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "webpack": "^5.0.0",
    "jest": "^29.0.0"
  }
}`,
		"src/index.ts": `import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { OTLPTraceExporter } from '@opentelemetry/exporter-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { UserInteractionInstrumentation } from '@opentelemetry/instrumentation-user-interaction';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import * as rrweb from 'rrweb';

export interface FrontendConfig {
  controlPlaneUrl: string;
  authKey: string;
  serviceName: string;
  environment?: string;
  enableSessionReplay?: boolean;
  redactionRules?: string[];
}

export class FrontendSDK {
  private sessionReplay?: rrweb.record;
  private stopRecording?: () => void;

  init(config: FrontendConfig): void {
    const resource = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: config.serviceName,
      [SemanticResourceAttributes.SERVICE_VERSION]: '1.0.0',
      [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: config.environment || 'production',
    });

    const exporter = new OTLPTraceExporter({
      url: \`\${config.controlPlaneUrl}/v1/traces\`,
      headers: {
        'Authorization': \`Bearer \${config.authKey}\`,
      },
    });

    const provider = new WebTracerProvider({
      resource,
      traceExporter: exporter,
    });

    provider.register();

    registerInstrumentations({
      instrumentations: [
        new DocumentLoadInstrumentation(),
        new FetchInstrumentation(),
        new UserInteractionInstrumentation(),
      ],
    });

    // Session replay
    if (config.enableSessionReplay) {
      this.startSessionReplay(config);
    }
  }

  private startSessionReplay(config: FrontendConfig): void {
    const events: rrweb.eventWithTime[] = [];
    
    this.sessionReplay = rrweb.record({
      emit(event) {
        events.push(event);
        // Batch and send events
        if (events.length >= 100) {
          this.sendEvents(events.splice(0, 100), config);
        }
      },
      maskAllInputs: true,
      maskTextSelector: config.redactionRules || [],
    });

    // Send remaining events on page unload
    window.addEventListener('beforeunload', () => {
      if (events.length > 0) {
        this.sendEvents(events, config);
      }
    });
  }

  private sendEvents(events: rrweb.eventWithTime[], config: FrontendConfig): void {
    fetch(\`\${config.controlPlaneUrl}/v1/sessions\`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': \`Bearer \${config.authKey}\`,
      },
      body: JSON.stringify({ events }),
    }).catch(console.error);
  }

  shutdown(): void {
    if (this.stopRecording) {
      this.stopRecording();
    }
  }
}

export default FrontendSDK;
`,
		"README.md": `# Frontend Observability SDK

Security observability SDK for web frontends (React, Vue, Angular, vanilla JS).

## Installation

\`\`\`bash
npm install @observability-platform/frontend-sdk
\`\`\`

## Usage

\`\`\`typescript
import { FrontendSDK } from '@observability-platform/frontend-sdk';

const sdk = new FrontendSDK();
sdk.init({
  controlPlaneUrl: 'https://control.example.com',
  authKey: process.env.AUTH_KEY,
  serviceName: 'my-web-app',
  environment: 'production',
  enableSessionReplay: true,
  redactionRules: ['[data-sensitive]', '.credit-card'],
});
\`\`\`
`,
	}

	for path, content := range files {
		if err := s.writeFile(path, content); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scaffolder) scaffoldJava() error {
	return s.writeFile("pom.xml", `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.observability</groupId>
    <artifactId>observability-agent-java</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>
    
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
`)
}

func (s *Scaffolder) scaffoldDotNet() error {
	return s.writeFile("ObservabilityAgent.csproj", `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <Version>1.0.0</Version>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="OpenTelemetry" Version="1.7.0" />
    <PackageReference Include="OpenTelemetry.Exporter.Otlp" Version="1.7.0" />
  </ItemGroup>
</Project>
`)
}

func (s *Scaffolder) scaffoldRuby() error {
	return s.writeFile("observability_agent.gemspec", `Gem::Specification.new do |spec|
  spec.name          = "observability-agent-ruby"
  spec.version       = "1.0.0"
  spec.summary       = "Security Observability Agent for Ruby"
  spec.files         = Dir["lib/**/*"]
  spec.require_paths = ["lib"]
end
`)
}

func (s *Scaffolder) scaffoldIOS() error {
	return s.writeFile("ObservabilityAgent.podspec", `Pod::Spec.new do |spec|
  spec.name         = "ObservabilityAgent"
  spec.version      = "1.0.0"
  spec.summary      = "Security Observability SDK for iOS"
  spec.source       = { :git => "", :tag => "#{spec.version}" }
  spec.source_files = "Sources/**/*.swift"
end
`)
}

func (s *Scaffolder) scaffoldAndroid() error {
	return s.writeFile("build.gradle", `plugins {
    id 'com.android.library'
    id 'org.jetbrains.kotlin.android'
}

android {
    compileSdk 34
    namespace 'com.observability.agent'
}
`)
}

func (s *Scaffolder) scaffoldDesktopWindows() error {
	return s.writeFile("AgentService.cs", `using System;
using System.ServiceProcess;

namespace ObservabilityAgent
{
    public partial class AgentService : ServiceBase
    {
        public AgentService()
        {
            InitializeComponent();
        }

        protected override void OnStart(string[] args)
        {
            // Initialize agent
        }

        protected override void OnStop()
        {
            // Shutdown agent
        }
    }
}
`)
}

func (s *Scaffolder) scaffoldDesktopMacOS() error {
	return s.writeFile("AgentDaemon.swift", `import Foundation

class AgentDaemon {
    func start() {
        // Initialize agent
    }
    
    func stop() {
        // Shutdown agent
    }
}
`)
}

func (s *Scaffolder) scaffoldDesktopLinux() error {
	return s.writeFile("agentd.go", `package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize agent daemon
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigChan
	log.Println("Shutting down...")
}
`)
}
