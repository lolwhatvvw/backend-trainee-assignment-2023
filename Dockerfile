FROM golang:1.21-alpine as build

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download

RUN go build -o main ./cmd/segment-api

CMD ["/app/main"]