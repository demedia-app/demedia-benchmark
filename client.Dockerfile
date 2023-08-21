FROM golang:1.20 as builder

WORKDIR /go/src/app
COPY ./ .

RUN go mod tidy

RUN cd client && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demedia-client .

FROM alpine:3.14

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/app/client/demedia-client /usr/local/bin/

ENTRYPOINT ["demedia-client"]
