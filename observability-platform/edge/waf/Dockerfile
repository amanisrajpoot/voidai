FROM nginx:alpine

# Install ModSecurity and OWASP CRS
RUN apk add --no-cache \
    modsecurity \
    modsecurity-nginx \
    wget \
    && mkdir -p /etc/nginx/modsec

# Download OWASP CRS
RUN wget -O /tmp/crs.tar.gz https://github.com/coreruleset/coreruleset/archive/v3.3.5.tar.gz \
    && tar -xzf /tmp/crs.tar.gz -C /etc/nginx/modsec \
    && mv /etc/nginx/modsec/coreruleset-3.3.5 /etc/nginx/modsec/coreruleset

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
COPY modsecurity.conf /etc/nginx/modsec/modsecurity.conf

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
