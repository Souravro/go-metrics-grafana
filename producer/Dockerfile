# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder
ENV APP_HOME /go/src/go-metrics-grafana/producer
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy producer directory
COPY . .

RUN go build -o /producer_service


FROM golang:1.19-alpine
ENV APP_HOME /go/src/go-metrics-grafana/producer
ENV GOPATH /go

WORKDIR "$APP_HOME"

# copy consumer config directory
COPY ./config/* ./config/

COPY --from=builder /producer_service .

ENTRYPOINT ["./producer_service"]
