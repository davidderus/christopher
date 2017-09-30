FROM golang:alpine

ENV SRC_DIR /go/src/github.com/davidderus/christopher

# Adding christopher sources
ADD . $SRC_DIR

RUN go install github.com/davidderus/christopher

# Exposing christopher
EXPOSE 8080
ENTRYPOINT ["/go/bin/christopher", "-c", "/christopher/config.toml"]

CMD ["--help"]
