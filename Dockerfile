FROM golang:1.27.0 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o kivo-server ./cmd/server && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o kivo-cli ./cmd/cli

FROM alpine:3.24 AS run

RUN apk add --no-cache ca-certificates ncurses-terminfo && \
    addgroup -S kivo && adduser -S kivo -G kivo

COPY --from=builder /app/kivo-server /usr/local/bin/kivo-server
COPY --from=builder /app/kivo-cli /usr/local/bin/kivo-cli

USER kivo

EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/kivo-server"]
