# Webcheckd: a website status check daemon

Webcheckd is a very simple utility to check the status of websites.
It runs in a loop and checks every 5 minutes if the site is up, expecting 
a **200 OK** http status.

It can be built as a standalone utility or executed as a container in Docker,
Podman or Kubernetes.

If the check fails it sends a notification using a user defined email account
to a user defined recipients list.

### Binary build
```
$ GO111MODULE=on go build
```

### Docker build
```
cd webcheckd
$ docker build -t webcheckd .
```

### Binary execution
```
$ ./webcheckd -url https://www.example.com \
  -from myaccount@gmail.com \
  -password mypassword \
  -to recipient1@gmail.com 
  -to recipient2@gmail.com
```

### Running as a container
```
$ docker run -e URL="https://www.example.com" \
  -e FROM=myaccount@gmail.com \
  -e PASSWORD=mypassword \
  -e TO=recipient1@gmail.com 
  webcheckd
```

Currently the docker execution supports only one recipient.
