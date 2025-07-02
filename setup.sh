#!/bin/bash

DOMAIN="tryit.selnastol.ru"
EMAIL="adagamov05@mail.ru"
DATA_PATH="/opt/certbot"

# Root check
if [ "$EUID" -ne 0 ]; then echo "Please run $0 as root." && exit; fi

# Create ALL required directories with proper permissions
echo "Setting up directory structure..."
mkdir -p "$DATA_PATH/www/.well-known/acme-challenge"
mkdir -p "$DATA_PATH/conf/live/$DOMAIN"
chmod -R 775 "$DATA_PATH"

# Create dummy certificates for each domain
echo "Creating dummy certificates..."
if [ ! -f "$DATA_PATH/conf/live/$DOMAIN/fullchain.pem" ]; then
openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
    -keyout "$DATA_PATH/conf/live/$DOMAIN/privkey.pem" \
    -out "$DATA_PATH/conf/live/$DOMAIN/fullchain.pem" \
    -subj "/CN=localhost"
fi

# Create test file for verification
echo "Creating test challenge file..."
echo "ACME_TEST" > "$DATA_PATH/www/.well-known/acme-challenge/test.txt"

# Start NGINX with dummy certs
echo "Starting NGINX with dummy certs..."
docker compose -f docker-compose.prod.yaml --env-file .env up -d frontend

# Verify the challenge file is accessible
echo "Verifying challenge endpoint..."
sleep 5
curl -v http://localhost/.well-known/acme-challenge/test.txt || \
  { echo "Challenge file not accessible!"; exit 1; }

# Request real certificates
echo "Requesting Let's Encrypt certificates for: ${DOMAIN}"
rm -rf "$DATA_PATH/conf/live/$DOMAIN"
rm -rf "$DATA_PATH/conf/archive/$DOMAIN"
rm -f    "$DATA_PATH/conf/renewal/$DOMAIN.conf"

docker compose -f docker-compose.prod.yaml run --rm --entrypoint "certbot certonly --webroot -w /var/www/certbot \
-d $DOMAIN \
--email $EMAIL \
--cert-name $DOMAIN \
--rsa-key-size 4096 \
--agree-tos \
--non-interactive \
--force-renewal" certbot

# Reload NGINX to pick up real certs
echo "Restarting NGINX with real certificates..."
docker compose -f docker-compose.prod.yaml restart frontend

echo "Done! Certificates issued for: ${DOMAIN}"