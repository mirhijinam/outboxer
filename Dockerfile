FROM golang:1.22-alpine

WORKDIR /go/src/outboxer

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/outboxer ./cmd/outboxer/main.go

EXPOSE ${HTTP_PORT}

CMD ["./bin/outboxer"]