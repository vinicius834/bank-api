FROM golang:1.12.17-buster

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o init .
CMD ["./init"]

