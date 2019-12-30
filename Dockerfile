FROM golang:latest
LABEL maintainer="Giovan Battista Salinetti <gbsalinetti@gmail.com>"

# Define workdir
WORKDIR /webcheckd

# Copy appliation files and dependencies
COPY go.mod go.sum main.go ./

# Downlaod modules and build
RUN go mod download && \
    go build

RUN mv webcheckd /usr/local/bin

EXPOSE 8080

USER 1001

ENTRYPOINT ["/usr/local/bin/webcheckd"]
CMD ["-h"]
