FROM golang:alpine

WORKDIR /
COPY oauth2_proxy /go/bin/oauth2_proxy

CMD oauth2_proxy  -policies-file  /config/policies  -provider oidc -client-id nginx-oauth2-proxy -client-secret ${CLIENT_SECRET} -cookie-secret ${CLIENT_SECRET} -oidc-issuer-url ${DEX_URL} -set-xauthrequest -pass-access-token -http-address ${LISTEN_IP}
