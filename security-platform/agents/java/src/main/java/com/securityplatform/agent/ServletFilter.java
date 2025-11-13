package com.securityplatform.agent;

import javax.servlet.*;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

/**
 * Servlet filter for automatic instrumentation.
 */
public class ServletFilter implements Filter {
    private Agent agent;

    public ServletFilter(Agent agent) {
        this.agent = agent;
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        // Initialization
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

        // Create span for request
        var span = agent.getTracer()
            .spanBuilder("http.request")
            .setAttribute("http.method", httpRequest.getMethod())
            .setAttribute("http.url", httpRequest.getRequestURL().toString())
            .startSpan();

        try (var scope = span.makeCurrent()) {
            chain.doFilter(request, response);
            span.setAttribute("http.status_code", httpResponse.getStatus());
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
        // Cleanup
    }
}
