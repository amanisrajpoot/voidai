package com.securityplatform.agent;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.InputStream;

/**
 * Factory for creating Agent instances from configuration.
 */
public class AgentFactory {
    private static final Logger logger = LoggerFactory.getLogger(AgentFactory.class);
    private static final ObjectMapper yamlMapper = new ObjectMapper(new YAMLFactory());
    private static final ObjectMapper jsonMapper = new ObjectMapper();

    /**
     * Create agent from YAML configuration file.
     */
    public static Agent fromYamlFile(String configPath) {
        try {
            AgentConfig config = yamlMapper.readValue(new File(configPath), AgentConfig.class);
            return new Agent(config);
        } catch (Exception e) {
            logger.error("Failed to load config from file: " + configPath, e);
            throw new RuntimeException("Failed to load agent configuration", e);
        }
    }

    /**
     * Create agent from YAML configuration stream.
     */
    public static Agent fromYamlStream(InputStream configStream) {
        try {
            AgentConfig config = yamlMapper.readValue(configStream, AgentConfig.class);
            return new Agent(config);
        } catch (Exception e) {
            logger.error("Failed to load config from stream", e);
            throw new RuntimeException("Failed to load agent configuration", e);
        }
    }

    /**
     * Create agent from JSON configuration file.
     */
    public static Agent fromJsonFile(String configPath) {
        try {
            AgentConfig config = jsonMapper.readValue(new File(configPath), AgentConfig.class);
            return new Agent(config);
        } catch (Exception e) {
            logger.error("Failed to load config from file: " + configPath, e);
            throw new RuntimeException("Failed to load agent configuration", e);
        }
    }

    /**
     * Create agent from configuration object.
     */
    public static Agent fromConfig(AgentConfig config) {
        return new Agent(config);
    }

    /**
     * Create agent with default configuration.
     */
    public static Agent createDefault() {
        AgentConfig config = new AgentConfig();
        config.setControlPlaneUrl(System.getenv().getOrDefault("SECURITY_PLATFORM_URL", "https://api.securityplatform.com"));
        config.setAuthKey(System.getenv("SECURITY_PLATFORM_AUTH_KEY"));
        config.setServiceName(System.getenv().getOrDefault("SERVICE_NAME", "java-service"));
        config.setEnvironment(System.getenv().getOrDefault("ENVIRONMENT", "production"));
        return new Agent(config);
    }
}
