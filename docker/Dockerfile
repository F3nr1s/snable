FROM golang:alpine AS builder
RUN apk --no-cache add ca-certificates build-base opusfile-dev
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
#RUN go test ./...
RUN go build  -o app ./*.go
#-ldflags="-extldflags=-static"
FROM alpine
RUN apk add --no-cache opusfile 
RUN mkdir /etc/snable
COPY --from=builder /build/app /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app", "-config", "/etc/snable/config.yml"]
