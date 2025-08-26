#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY *.go .
COPY *.mod .
COPY *.yaml .
RUN go mod download
COPY . .
COPY config.yaml .
RUN go get -d -v ./.
RUN go build -o /go/bin/app -v ./.

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/src/app/config.yaml .

ENTRYPOINT /app
LABEL Name=go-starter-auth Version=0.0.1
EXPOSE 9922
CMD ["/app"]
