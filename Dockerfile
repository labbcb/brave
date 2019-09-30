FROM golang:1.13 AS builder
WORKDIR /go/src/github.com/labbcb/brave
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o brave .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/labbcb/brave/brave /usr/local/bin/brave
ENTRYPOINT ["geni"]