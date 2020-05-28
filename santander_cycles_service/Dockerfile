FROM golang:1.14
WORKDIR /go/src
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/app .

FROM alpine:latest
# RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/app /bin
CMD ["app"]
