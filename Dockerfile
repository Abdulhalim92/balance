FROM golang:latest

WORKDIR /balance

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main

ENV PORT 8080
EXPOSE 8080

CMD ["/balance/main"]