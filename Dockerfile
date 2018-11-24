FROM alpine:latest
MAINTAINER iYOCHU Nigeria Ltd

RUN apk add --no-cache ca-certificates

# Add the executable
COPY buyworker /buyworker



CMD ["/buyworker"]