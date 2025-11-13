package com.securityplatform.agent;

import java.util.List;
import java.util.Map;

/**
 * Configuration for the Security Platform Agent.
 */
public class AgentConfig {
    private String controlPlaneUrl;
    private String authKey;
    private String serviceName;
    private String environment;
    private String namespace;
    private String otlpEndpoint;
    private TelemetryConfig telemetry;
    private List<RedactionRule> redactionRules;
    private PolicyConfig policy;
    private SecurityConfig security;

    public AgentConfig() {
        this.telemetry = new TelemetryConfig();
        this.policy = new PolicyConfig();
        this.security = new SecurityConfig();
    }

    // Getters and setters
    public String getControlPlaneUrl() { return controlPlaneUrl; }
    public void setControlPlaneUrl(String controlPlaneUrl) { this.controlPlaneUrl = controlPlaneUrl; }

    public String getAuthKey() { return authKey; }
    public void setAuthKey(String authKey) { this.authKey = authKey; }

    public String getServiceName() { return serviceName; }
    public void setServiceName(String serviceName) { this.serviceName = serviceName; }

    public String getEnvironment() { return environment; }
    public void setEnvironment(String environment) { this.environment = environment; }

    public String getNamespace() { return namespace; }
    public void setNamespace(String namespace) { this.namespace = namespace; }

    public String getOtlpEndpoint() { return otlpEndpoint; }
    public void setOtlpEndpoint(String otlpEndpoint) { this.otlpEndpoint = otlpEndpoint; }

    public TelemetryConfig getTelemetry() { return telemetry; }
    public void setTelemetry(TelemetryConfig telemetry) { this.telemetry = telemetry; }

    public List<RedactionRule> getRedactionRules() { return redactionRules; }
    public void setRedactionRules(List<RedactionRule> redactionRules) { this.redactionRules = redactionRules; }

    public PolicyConfig getPolicy() { return policy; }
    public void setPolicy(PolicyConfig policy) { this.policy = policy; }

    public SecurityConfig getSecurity() { return security; }
    public void setSecurity(SecurityConfig security) { this.security = security; }

    public static class TelemetryConfig {
        private int batchSize = 100;
        private String batchTimeout = "5s";
        private String exportTimeout = "30s";
        private int maxQueueSize = 2048;

        public int getBatchSize() { return batchSize; }
        public void setBatchSize(int batchSize) { this.batchSize = batchSize; }

        public String getBatchTimeout() { return batchTimeout; }
        public void setBatchTimeout(String batchTimeout) { this.batchTimeout = batchTimeout; }

        public String getExportTimeout() { return exportTimeout; }
        public void setExportTimeout(String exportTimeout) { this.exportTimeout = exportTimeout; }

        public int getMaxQueueSize() { return maxQueueSize; }
        public void setMaxQueueSize(int maxQueueSize) { this.maxQueueSize = maxQueueSize; }
    }

    public static class PolicyConfig {
        private String mode = "observe"; // "observe" or "block"
        private boolean autoEnableBlocking = false;
        private int observePeriodHours = 48;
        private List<Map<String, Object>> rules;

        public String getMode() { return mode; }
        public void setMode(String mode) { this.mode = mode; }

        public boolean isAutoEnableBlocking() { return autoEnableBlocking; }
        public void setAutoEnableBlocking(boolean autoEnableBlocking) { this.autoEnableBlocking = autoEnableBlocking; }

        public int getObservePeriodHours() { return observePeriodHours; }
        public void setObservePeriodHours(int observePeriodHours) { this.observePeriodHours = observePeriodHours; }

        public List<Map<String, Object>> getRules() { return rules; }
        public void setRules(List<Map<String, Object>> rules) { this.rules = rules; }
    }

    public static class SecurityConfig {
        private boolean mtlsEnabled = true;
        private String certificatePath;
        private String keyPath;
        private String caBundlePath;

        public boolean isMtlsEnabled() { return mtlsEnabled; }
        public void setMtlsEnabled(boolean mtlsEnabled) { this.mtlsEnabled = mtlsEnabled; }

        public String getCertificatePath() { return certificatePath; }
        public void setCertificatePath(String certificatePath) { this.certificatePath = certificatePath; }

        public String getKeyPath() { return keyPath; }
        public void setKeyPath(String keyPath) { this.keyPath = keyPath; }

        public String getCaBundlePath() { return caBundlePath; }
        public void setCaBundlePath(String caBundlePath) { this.caBundlePath = caBundlePath; }
    }
}
