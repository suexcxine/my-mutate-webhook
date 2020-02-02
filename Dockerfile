FROM alpine:latest
RUN apk --no-cache add ca-certificates && mkdir -p /app
WORKDIR /app
ADD my-mutate .
ADD ssl ssl
ENTRYPOINT ["/app/my-mutate"]

