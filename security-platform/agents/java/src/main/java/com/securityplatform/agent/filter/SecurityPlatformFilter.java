package com.securityplatform.agent.filter;

import com.securityplatform.agent.Agent;
import com.securityplatform.agent.PolicyResult;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.servlet.*;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

/**
 * Servlet filter for Security Platform Agent.
 */
public class SecurityPlatformFilter implements Filter {
    private static final Logger logger = LoggerFactory.getLogger(SecurityPlatformFilter.class);
    
    private Agent agent;

    public SecurityPlatformFilter(Agent agent) {
        this.agent = agent;
    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {
        if (!(request instanceof HttpServletRequest) || !(response instanceof HttpServletResponse)) {
            chain.doFilter(request, response);
            return;
        }

        HttpServletRequest httpRequest = (HttpServletRequest) request;
        HttpServletResponse httpResponse = (HttpServletResponse) response;

        // Check policy
        Map<String, Object> context = new HashMap<>();
        context.put("method", httpRequest.getMethod());
        context.put("path", httpRequest.getRequestURI());
        context.put("remote_addr", httpRequest.getRemoteAddr());
        PolicyResult policyResult = agent.checkPolicy("http_request", context);

        if (!policyResult.isAllowed()) {
            logger.warn("Request blocked by policy: {}", policyResult.getReason());
            httpResponse.setStatus(403);
            httpResponse.getWriter().write("Request blocked by security policy");
            return;
        }

        // Continue with request
        chain.doFilter(request, response);
    }
}
