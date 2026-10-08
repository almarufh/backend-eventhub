FROM golang:1.27.1-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o backend ./cmd/main.go

FROM alpine:3.24.2

WORKDIR /server

COPY --from=builder /app/backend /server/backend

CMD [ "/server/backend" ] 