#!/bin/sh

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout charts/code-editor/tls.key -out charts/code-editor/tls.crt -subj "/C=XX/ST=Italy/L=Empoli/O=SUSE/OU=ECM/CN=code-editor" 2>/dev/null

echo "Self-signed certificate successfully generated"