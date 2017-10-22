#!/usr/bin/env bash

docker rm -f simplifyServer
docker pull ssharif6/simplify
docker run -d \
-p 443:443 \
--name simplifyServer \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=/etc/letsencrypt/live/simplify.api.shaheensharifian.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/simplify.api.shaheensharifian.me/privkey.pem \
-e GOOGLE_APPLICATION_CREDENTIALS=/Simplify-a6b380e0b81c.json \
ssharif6/simplify
