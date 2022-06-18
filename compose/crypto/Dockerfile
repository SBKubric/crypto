FROM golang:1.17-alpine as build-stage

RUN mkdir -p /app

WORKDIR /app

COPY . /app
RUN go mod download

RUN go build -o crypto main.go

FROM alpine:latest

WORKDIR /

COPY --from=build-stage /app/crypto .

EXPOSE 9999

ENTRYPOINT [ "/crypto" ]
