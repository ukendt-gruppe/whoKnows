#!/bin/bash
if [ ! -f /etc/letsencrypt/live/monkbusiness.dk/fullchain.pem ]; then
    certbot --nginx \
            --non-interactive \
            --agree-tos \
            --email sifo0001@stud.kea.dk \
            -d monkbusiness.dk \
            --redirect
fi 