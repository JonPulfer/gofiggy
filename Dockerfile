FROM golang:1.14 AS builder
LABEL stage=builder

WORKDIR /app/

COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./gofiggy  ./cmd/gofiggy

FROM alpine
RUN apk add --no-cache libc6-compat ca-certificates apache2-utils
COPY --from=builder /app/gofiggy ./
ENTRYPOINT ["./gofiggy"]
