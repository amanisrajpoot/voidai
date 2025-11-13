FROM nginx:alpine

# Install ModSecurity and dependencies
RUN apk add --no-cache \
    modsecurity \
    modsecurity-nginx \
    wget \
    curl

# Copy ModSecurity configuration
COPY modsecurity.conf /etc/nginx/modsecurity.conf
COPY nginx.conf /etc/nginx/nginx.conf
COPY crs-setup.conf /etc/nginx/crs-setup.conf

# Download OWASP CRS rules
RUN mkdir -p /etc/nginx/crs && \
    wget -q -O - https://github.com/coreruleset/coreruleset/archive/v3.3.5.tar.gz | \
    tar -xz -C /etc/nginx/crs --strip-components=1

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
