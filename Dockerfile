FROM docker.io/library/golang:1.24.2-bookworm AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 go build -trimpath -o goproxygo cmd/goproxygo/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/goproxygo /goproxygo

EXPOSE 8080
ENTRYPOINT ["/goproxygo", "-host", "0.0.0.0"]

