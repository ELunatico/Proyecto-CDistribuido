FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o balanceador main.go

EXPOSE 9000

CMD ["./balanceador"]