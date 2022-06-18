FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./


RUN go build -o crypto main.go

EXPOSE 9999

CMD [ "/main" ]
