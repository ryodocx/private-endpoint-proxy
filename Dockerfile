FROM golang:1.20.2-alpine
WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go install

FROM alpine:3.17.2
COPY --from=0 /go/bin/private-endpoint-proxy /usr/local/bin/
COPY ./example/config.yaml /var/config.yaml
CMD [ "private-endpoint-proxy", "-c", "/var/config.yaml" ]
