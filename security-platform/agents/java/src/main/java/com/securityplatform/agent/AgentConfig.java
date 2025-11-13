package com.securityplatform.agent;

import java.util.List;
import java.util.Map;

/**
 * Configuration for Security Platform Agent
 */
public class AgentConfig {
    private String controlPlaneUrl;
    private String authKey;
    private String serviceName;
    private String version;
    private String environment;
    private String otlpEndpoint;
    private List<String> redactionRules;
    private String localPolicy; // "observe" or "block"
    private Integer telemetryBatchSize;
    private Map<String, String> customAttributes;

    public String getControlPlaneUrl() {
        return controlPlaneUrl;
    }

    public void setControlPlaneUrl(String controlPlaneUrl) {
        this.controlPlaneUrl = controlPlaneUrl;
    }

    public String getAuthKey() {
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

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
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

    public List<String> getRedactionRules() {
        return redactionRules;
    }

    public void setRedactionRules(List<String> redactionRules) {
        this.redactionRules = redactionRules;
    }

    public String getLocalPolicy() {
        return localPolicy;
    }

    public void setLocalPolicy(String localPolicy) {
        this.localPolicy = localPolicy;
    }

    public Integer getTelemetryBatchSize() {
        return telemetryBatchSize;
    }

    public void setTelemetryBatchSize(Integer telemetryBatchSize) {
        this.telemetryBatchSize = telemetryBatchSize;
    }

    public Map<String, String> getCustomAttributes() {
        return customAttributes;
    }

    public void setCustomAttributes(Map<String, String> customAttributes) {
        this.customAttributes = customAttributes;
    }
}
