FROM ubuntu:latest

WORKDIR /app
COPY ./cmd/accrual/bbhelper ./bbhelper

CMD ["./bbhelper",  "-r", "http://accrual:8080"]
