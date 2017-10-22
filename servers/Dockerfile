FROM alpine
RUN apk add --no-cache ca-certificates
COPY servers /servers
EXPOSE 443
ENTRYPOINT ["/servers"]
