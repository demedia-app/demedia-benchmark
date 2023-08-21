FROM golang:1.20 as builder

WORKDIR /go/src/app
COPY ./ .

RUN go mod tidy

RUN cd peer && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demedia-peer .

FROM alpine:3.14

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/app/peer/demedia-peer /usr/local/bin/

ENTRYPOINT ["demedia-peer"]
