# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder
ENV APP_HOME /go/src/go-metrics-grafana/dashboard
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy dashboard directory
COPY . .

RUN go build -o /dashboard_service


FROM golang:1.19-alpine
ENV APP_HOME /go/src/go-metrics-grafana/dashboard
ENV GOPATH /go

WORKDIR "$APP_HOME"

# copy dashboard config directory
COPY ./config/* ./config/

COPY --from=builder /dashboard_service .

EXPOSE 2121

ENTRYPOINT ["./dashboard_service"]
