FROM nginx:latest

# Install certbot and its nginx plugin
RUN apt-get update && \
    apt-get install -y certbot python3-certbot-nginx && \
    rm -rf /var/lib/apt/lists/*

# Create required directories
RUN mkdir -p /var/www/certbot /etc/nginx/conf.d

# Copy nginx configurations
COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./nginx-init.conf /etc/nginx/conf.d/default.conf
COPY ./nginx-ssl.conf /etc/nginx/conf.d/ssl.conf.disabled

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]