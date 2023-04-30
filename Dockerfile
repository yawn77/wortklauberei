FROM golang:1.20-alpine as builder

WORKDIR /app
ARG VERSION

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
RUN go build -ldflags="-X main.version=$VERSION" -o wordklauberei ./cmd/wordklauberei/main.go

FROM alpine:latest

RUN apk update && \
    apk add --no-cache tzdata

COPY --from=builder /app/wordklauberei /app/wordklauberei

CMD [ "/app/wordklauberei" ]
