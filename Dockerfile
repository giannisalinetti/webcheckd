FROM golang:latest
LABEL maintainer="Giovan Battista Salinetti <gbsalinetti@gmail.com>"

ENV URL=https://www.paradewedding.com \
    FROM=sender@gmail.com \
    PASSWORD=mypassword \
    TO=recipient@gmail.com

# Define workdir
WORKDIR /webcheckd

# Copy appliation files and dependencies
COPY go.mod go.sum main.go ./

# Downlaod modules and build
RUN go mod download && \
    go build

RUN mv webcheckd /usr/local/bin

ENTRYPOINT ["sh", "-c", "/usr/local/bin/webcheckd"]
CMD -url $URL -from $FROM -password $PASSWORD -to $TO
