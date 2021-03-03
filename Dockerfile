# Create Builder & Build
FROM golang:1.16.0-alpine3.13 as builder

RUN apk update && apk add --no-cache git && apk add ca-certificates
RUN adduser -D -g '' appuser

# Create Container
FROM scratch 

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY  k8s-global-objects /app

USER appuser
ENTRYPOINT ["/app"]
