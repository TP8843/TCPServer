FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM alpine:3.21.3

WORKDIR /

COPY --from=builder /server /server

EXPOSE 3000

CMD ["/server"]