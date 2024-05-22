#!/bin/sh

HOST_DOMAIN=$(yq '.domain' helm-charts/code-editor/values.yaml)

echo "Certificate CN: '/CN=${HOST_DOMAIN}'"

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout helm-charts/code-editor/assets/tls.key -out helm-charts/code-editor/assets/tls.crt -subj "/CN=${HOST_DOMAIN}" 2>/dev/null

echo "Self-signed certificate successfully generated"