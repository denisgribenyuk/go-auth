ARG GO_VERSION="1.19"
ARG ALPINE_VERSION="3.16"
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /
COPY . .

RUN apk update && \
    env GOWORK=off go build -o back_auth

FROM alpine:${ALPINE_VERSION}
WORKDIR /app
COPY --from=builder /back_auth /app/

EXPOSE 8080
CMD ["./back_auth"]