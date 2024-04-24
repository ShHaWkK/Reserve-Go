FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o reservationservice .

FROM scratch
COPY --from=builder /app/reservationservice /reservationservice

FROM mysql:8.0
COPY BDD.sql /docker-entrypoint-initdb.d/

CMD ["/reservationservice"]
