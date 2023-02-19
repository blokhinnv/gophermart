FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY ./internal  ./internal
COPY ./cmd  ./cmd

RUN go build -o ./gophermart ./cmd/gophermart/main.go
CMD ["./gophermart"]
