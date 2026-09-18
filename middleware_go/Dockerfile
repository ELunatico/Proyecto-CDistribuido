FROM golang:alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o orquestador main.go

EXPOSE 8080
CMD ["./orquestador"]