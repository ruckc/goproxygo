FROM docker.io/library/golang:1.17.7-bullseye AS build

WORKDIR /app

COPY . /app

RUN go build -o goproxygo cmd/goproxygo/main.go

FROM docker.io/library/debian:stable-slim

COPY --from=build /app/goproxygo /goproxygo

EXPOSE 8080
ENTRYPOINT ["/goproxygo", "-host", "0.0.0.0"]

