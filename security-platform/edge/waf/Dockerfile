FROM nginx:alpine

# Install ModSecurity and dependencies
RUN apk add --no-cache \
    modsecurity \
    modsecurity-nginx \
    wget \
    curl

# Install OWASP Core Rule Set
RUN mkdir -p /etc/nginx/modsec && \
    wget -O /tmp/owasp.tar.gz https://github.com/coreruleset/coreruleset/archive/v3.3.5.tar.gz && \
    tar -xzf /tmp/owasp.tar.gz -C /etc/nginx/modsec --strip-components=1 && \
    rm /tmp/owasp.tar.gz

# Copy ModSecurity configuration
COPY modsecurity.conf /etc/nginx/modsec/
COPY crs-setup.conf /etc/nginx/modsec/coreruleset/

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf

# Copy security platform WAF plugin
COPY waf-plugin.lua /etc/nginx/lua/

# Install Lua dependencies (if using lua-resty-http)
RUN apk add --no-cache lua5.1 lua5.1-dev luarocks && \
    luarocks install lua-resty-http

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
