package com.securityplatform.agent;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;
import java.util.Map;

/**
 * Configuration for the Security Platform Agent.
 */
public class AgentConfig {
    @JsonProperty("control_plane_url")
    private String controlPlaneUrl = "https://api.securityplatform.com";

    @JsonProperty("auth_key")
    private String authKey;

    @JsonProperty("service_name")
    private String serviceName = "java-service";

    @JsonProperty("environment")
    private String environment = "production";

    @JsonProperty("otlp_endpoint")
    private String otlpEndpoint;

    @JsonProperty("telemetry")
    private TelemetryConfig telemetry = new TelemetryConfig();

    @JsonProperty("redaction_rules")
    private List<String> redactionRules;

    @JsonProperty("policy")
    private PolicyConfig policy = new PolicyConfig();

    // Getters and setters
    public String getControlPlaneUrl() {
        return controlPlaneUrl != null ? controlPlaneUrl : System.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL");
    }

    public void setControlPlaneUrl(String controlPlaneUrl) {
        this.controlPlaneUrl = controlPlaneUrl;
    }

    public String getAuthKey() {
        if (authKey == null || authKey.isEmpty()) {
            return System.getenv("SECURITY_PLATFORM_AUTH_KEY");
        }
        return authKey;
    }

    public void setAuthKey(String authKey) {
        this.authKey = authKey;
    }

    public String getServiceName() {
        return serviceName;
    }

    public void setServiceName(String serviceName) {
        this.serviceName = serviceName;
    }

    public String getEnvironment() {
        return environment;
    }

    public void setEnvironment(String environment) {
        this.environment = environment;
    }

    public String getOtlpEndpoint() {
        return otlpEndpoint;
    }

    public void setOtlpEndpoint(String otlpEndpoint) {
        this.otlpEndpoint = otlpEndpoint;
    }

    public TelemetryConfig getTelemetry() {
        return telemetry;
    }

    public void setTelemetry(TelemetryConfig telemetry) {
        this.telemetry = telemetry;
    }

    public List<String> getRedactionRules() {
        return redactionRules;
    }

    public void setRedactionRules(List<String> redactionRules) {
        this.redactionRules = redactionRules;
    }

    public PolicyConfig getPolicy() {
        return policy;
    }

    public void setPolicy(PolicyConfig policy) {
        this.policy = policy;
    }

    public static class TelemetryConfig {
        @JsonProperty("batch_size")
        private int batchSize = 100;

        @JsonProperty("batch_timeout")
        private int batchTimeout = 5; // seconds

        @JsonProperty("export_timeout")
        private int exportTimeout = 30; // seconds

        @JsonProperty("max_queue_size")
        private int maxQueueSize = 2048;

        public int getBatchSize() {
            return batchSize;
        }

        public void setBatchSize(int batchSize) {
            this.batchSize = batchSize;
        }

        public int getBatchTimeout() {
            return batchTimeout;
        }

        public void setBatchTimeout(int batchTimeout) {
            this.batchTimeout = batchTimeout;
        }

        public int getExportTimeout() {
            return exportTimeout;
        }

        public void setExportTimeout(int exportTimeout) {
            this.exportTimeout = exportTimeout;
        }

        public int getMaxQueueSize() {
            return maxQueueSize;
        }

        public void setMaxQueueSize(int maxQueueSize) {
            this.maxQueueSize = maxQueueSize;
        }
    }

    public static class PolicyConfig {
        @JsonProperty("mode")
        private String mode = "observe";

        public String getMode() {
            return mode;
        }

        public void setMode(String mode) {
            this.mode = mode;
        }
    }
}
