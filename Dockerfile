FROM golang:1.22
WORKDIR /go/src
COPY internal ./internal
COPY main.go .
COPY go.sum .
COPY go.mod .

RUN go build

ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["./user-service"]
