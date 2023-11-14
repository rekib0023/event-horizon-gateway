FROM golang:1.20.4 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o event-horizon-gateway

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/event-horizon-gateway .

ENTRYPOINT ["./event-horizon-gateway"]
