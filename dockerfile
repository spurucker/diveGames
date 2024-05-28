FROM golang:1.18-alpine

USER root

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod -R 755 /app

RUN go build -o main ./main

EXPOSE 8080

CMD ["go", "run", "./main"]