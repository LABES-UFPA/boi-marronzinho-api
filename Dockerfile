FROM golang:1.22.6-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -tags musl -o main

FROM alpine:3.18

WORKDIR /app

COPY --from=build /app/main .

CMD ["./main"]