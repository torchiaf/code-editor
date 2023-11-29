104.248.31.127

http://104.248.31.127/code-editor/mFYRM/?folder=/git/code-editor

https://localhost/code-editor/mFYRM/?folder=/git/code-editor

https://epinio.127.0.0.1.nip.io/code-editor/mFYRM/?folder=/git/code-editor

https://epinio.104.248.31.127.nip.io/code-editor/api/v1/auth
https://epinio.104.248.31.127.nip.io/code-editor/mFYRM/login?folder=/git/code-editor


Changes auth -> to login

ghcr.io/torchiaf/epinio-ui:sha256-8f99abf7036cea6635eafd7946f0449de67258f360b4306573cfea86afeec7c8.sig


openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ~/certs/MyKey.key -out ~/certs/MyCertificate.crt -subj "/C=XX/ST=Italy/L=Empoli/O=SUSE/OU=ECM/CN=pippo"




---- create secret with minio and use it for tls connection

openssl req -newkey rsa:2048 -nodes -keyout tls.key -x509 -days 3650 -out tls.crt
kubectl -n code-editor create secret tls minio-tls --key=tls.key --cert=tls.crt