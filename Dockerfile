FROM golang:latest AS build

WORKDIR /app

COPY *.go go.mod go.sum /app/

COPY dom /app/dom

RUN go mod tidy &&\
    go build -o go-sample-app .

FROM debian:stable-slim AS bin

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

EXPOSE 9090

COPY --from=build /app/go-sample-app /app/

CMD [ "/app/go-sample-app" ]
