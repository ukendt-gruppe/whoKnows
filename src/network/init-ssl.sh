#!/bin/bash

# Wait for nginx to start
sleep 5

# Try to get the certificate
certbot certonly --webroot -w /var/www/certbot --email sifo0001@stud.kea.dk -d monkbusiness.dk --agree-tos --non-interactive

# If successful, enable SSL configuration
if [ -f /etc/letsencrypt/live/monkbusiness.dk/fullchain.pem ]; then
    mv /etc/nginx/conf.d/ssl.conf.disabled /etc/nginx/conf.d/ssl.conf
    nginx -s reload
fi