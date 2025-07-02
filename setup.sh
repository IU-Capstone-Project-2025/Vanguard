#!/bin/bash

DOMAIN="tryit.selnastol.ru"
EMAIL="adagamov05@mail.ru"
DATA_PATH="/opt/certbot"

# Root check
if [ "$EUID" -ne 0 ]; then echo "Please run $0 as root." && exit; fi

# Create dummy certificates for each domain
echo "Creating dummy certificates..."
mkdir -p "$DATA_PATH/conf/live/$DOMAIN"
if [ ! -f "$DATA_PATH/conf/live/$DOMAIN/fullchain.pem" ]; then
openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
    -keyout "$DATA_PATH/conf/live/$DOMAIN/privkey.pem" \
    -out "$DATA_PATH/conf/live/$DOMAIN/fullchain.pem" \
    -subj "/CN=localhost"
fi

# Start NGINX with dummy certs
echo "Starting NGINX with dummy certs..."
docker compose up -d nginx

# Give NGINX a moment to spin up
sleep 5

# Request real certificates
echo "Requesting Let's Encrypt certificates for: ${DOMAIN}"
rm -rf "$DATA_PATH/conf/live/$DOMAIN"
rm -rf "$DATA_PATH/conf/archive/$DOMAIN"
rm -f    "$DATA_PATH/conf/renewal/$DOMAIN.conf"

docker compose run --rm --entrypoint "certbot certonly --webroot -w /var/www/certbot \
-d $DOMAIN \
--email $EMAIL \
--cert-name $DOMAIN \
--rsa-key-size 4096 \
--agree-tos \
--non-interactive \
--force-renewal" certbot

# Reload NGINX to pick up real certs
echo "Restarting NGINX with real certificates..."
docker compose restart nginx

echo "Done! Certificates issued for: ${DOMAIN}"