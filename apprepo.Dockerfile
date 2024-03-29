FROM golang:1.18.9 as builder
WORKDIR /app

RUN apt-get install git
RUN go install github.com/go-delve/delve/cmd/dlv@v1.9.1

COPY go.mod go.sum ./
COPY cmd/apprepository-controller/ cmd/apprepository-controller/

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o apprepository-controller cmd/apprepository-controller/main.go
RUN CGO_ENABLED=0 go build -a -installsuffix cgo ./cmd/apprepository-controller

#### BASE ####
FROM gcr.io/distroless/static-debian11:nonroot AS base

WORKDIR /app

COPY --from=builder /app/apprepository-controller /apprepository-controller
COPY --from=builder /go/bin/dlv /dlv

CMD ["/apprepository-controller"]
