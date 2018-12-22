FROM golang:latest AS Builder

COPY . /go/src/app
WORKDIR /go/src/app

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build

FROM alpine:latest AS Runner
RUN apk add --update ca-certificates
COPY  --from=Builder /go/src/app/app /usr/local/bin/app

ENTRYPOINT ["app"]

