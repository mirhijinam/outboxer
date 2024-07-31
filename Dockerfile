FROM --platform=linux/amd64 golang:1.22-alpine

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update && \
     apk add \
        alpine-sdk \
        linux-headers \
        gcompat \
        libstdc++ \
        gcc \
        g++


WORKDIR /go/src/outboxer

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -tags musl --ldflags "-extldflags -static" -o ./bin/outboxer ./cmd/outboxer/main.go

EXPOSE ${HTTP_PORT}

CMD ["./bin/outboxer"]