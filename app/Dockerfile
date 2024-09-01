FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o nous ./cmd/api

EXPOSE 8080

RUN mkdir -p /app/data && chmod 777 /app/data

CMD ["./nous"]