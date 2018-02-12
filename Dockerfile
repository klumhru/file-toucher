FROM golang as builder

COPY . /go/src/app

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app /go/src/app/main.go

FROM alpine

COPY --from=builder /go/bin/app /usr/bin/app

ENTRYPOINT ["/usr/bin/app"]
