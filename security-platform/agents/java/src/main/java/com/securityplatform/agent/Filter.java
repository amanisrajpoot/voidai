package com.securityplatform.agent;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import java.io.IOException;
import java.util.Map;
import java.util.HashMap;

/**
 * Servlet Filter for Security Platform Agent
 */
public class Filter implements jakarta.servlet.Filter {
    private Agent agent;
    private Redaction redaction;
    private PolicyEngine policyEngine;

    public Filter(Agent agent) {
        this.agent = agent;
        this.redaction = new Redaction(agent.getConfig().getRedactionRules());
        this.policyEngine = new PolicyEngine(agent.getConfig());
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        // No initialization needed
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

        // Create span
        var span = agent.getTracer()
            .spanBuilder("http.request")
            .startSpan();

        try (var scope = span.makeCurrent()) {
            // Capture request data
            Map<String, Object> requestData = new HashMap<>();
            requestData.put("method", httpRequest.getMethod());
            requestData.put("uri", httpRequest.getRequestURI());
            requestData.put("query", httpRequest.getQueryString());
            
            Map<String, String> headers = new HashMap<>();
            httpRequest.getHeaderNames().asIterator()
                .forEachRemaining(name -> headers.put(name, httpRequest.getHeader(name)));
            requestData.put("headers", redaction.redact(headers));

            // Check policy
            PolicyResult policyResult = policyEngine.checkPolicy("http_request", requestData);
            
            if ("block".equals(agent.getConfig().getLocalPolicy()) && !policyResult.isAllowed()) {
                httpResponse.setStatus(403);
                httpResponse.getWriter().write("Request blocked by security policy");
                span.setAttribute("security.blocked", true);
                span.setAttribute("security.reason", policyResult.getReason());
                return;
            }

            // Continue with request
            chain.doFilter(request, response);

            // Capture response data
            span.setAttribute("http.status_code", httpResponse.getStatus());
            span.setAttribute("security.policy.allowed", policyResult.isAllowed());
            
            if (!policyResult.isAllowed()) {
                span.setAttribute("security.policy.reason", policyResult.getReason());
            }

        } catch (Exception e) {
            span.recordException(e);
            span.setAttribute("error", true);
            throw e;
        } finally {
            span.end();
        }
    }

    @Override
    public void destroy() {
        // No cleanup needed
    }
}
