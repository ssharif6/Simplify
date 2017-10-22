FROM alpine
RUN apk add --no-cache ca-certificates
COPY servers /servers
COPY Simplify-a6b380e0b81c.json Simplify-a6b380e0b81c.json
EXPOSE 443
ENTRYPOINT ["/servers"]
