FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod main.go ./
RUN go get github.com/useflagly/sdk-go@latest && go mod tidy && \
    GOFILE=$(find /go/pkg/mod/github.com/useflagly -name "client.go" | head -1) && \
    chmod +w "$GOFILE" && \
    sed -i 's/req.Header.Set("Authorization", "Bearer "+c.token)/req.Header.Set("apiKey", c.token)/' "$GOFILE"

RUN go build -o example .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/example .

CMD ["./example"]
