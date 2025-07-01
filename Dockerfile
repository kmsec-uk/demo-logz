FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
ADD https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css /app/static/pico.min.css
RUN go build -v -o /app/server

FROM alpine:3
RUN addgroup -g 101 -S web && adduser -h /app -u 1001 -D web -G web
COPY --from=builder /app/server /app/server
WORKDIR /app
USER web
CMD ["/app/server"]