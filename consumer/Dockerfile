# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder
ENV APP_HOME /go/src/go-metrics-grafana/consumer
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy consumer directory
COPY . .

RUN go build -o /consumer_service


FROM golang:1.19-alpine
ENV APP_HOME /go/src/go-metrics-grafana/consumer
ENV GOPATH /go

WORKDIR "$APP_HOME"

# copy consumer config directory
COPY ./config/* ./config/

COPY --from=builder /consumer_service .

EXPOSE 8080

ENTRYPOINT ["./consumer_service"]
