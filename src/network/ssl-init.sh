#!/bin/bash

# Create webroot directory
mkdir -p /var/www/certbot

# Wait for nginx to start
sleep 5

# Initialize SSL configuration file
SSL_CONF="/etc/nginx/conf.d/ssl.conf"
touch $SSL_CONF

if [ ! -f /etc/letsencrypt/live/monkbusiness.dk/fullchain.pem ]; then
    echo "No SSL certificate found. Obtaining one..."
    
    # Get the certificate
    certbot certonly --webroot \
            --webroot-path /var/www/certbot \
            --non-interactive \
            --agree-tos \
            --email monk@monkbusiness.dk \
            -d monkbusiness.dk

    # If certificate was obtained successfully, create SSL configuration
    if [ -f /etc/letsencrypt/live/monkbusiness.dk/fullchain.pem ]; then
        cat > $SSL_CONF <<EOF
server {
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name monkbusiness.dk;

    ssl_certificate /etc/letsencrypt/live/monkbusiness.dk/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/monkbusiness.dk/privkey.pem;

    location /static/ {
        alias /app/frontend/static/;
        expires 30d;
        add_header Cache-Control "public, no-transform";
    }

    location / {
        proxy_pass http://whoknows_go:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF
        # Update HTTP server to redirect to HTTPS
        sed -i 's/proxy_pass/return 301 https:\/\/$host$request_uri;#/' /etc/nginx/nginx.conf
    fi
fi

# Exit the script - let the main nginx process handle everything
exit 0