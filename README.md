# Webcheckd: a website status check daemon

Webcheckd is a very simple utility to check the status of websites.
It runs in a loop and checks every 5 minutes if the monitored site are up, 
expecting a **200 OK** http status in front of a simple HTTP GET request.

It can be built as a standalone utility or executed as a container in Docker,
Podman or Kubernetes.

If the check fails it sends a notification using a user defined email account
to a user defined recipients list.

### Binary build
To build the application use the **make** command along with the provided 
Makefile.
```
$ make build
```

### Docker build
The make command can also be used to build, tag and push the image.
```
$ make image; make tag; make push
```

### Binary execution
This is a basic example of execution. More than one url can be checked in 
parallel.
```
$ ./webcheckd \
  -url https://www.example.com \
  -url https://www.myapp.com \
  -from myaccount@gmail.com \
  -password mypassword \
  -to recipient1@gmail.com 
  -to recipient2@gmail.com
```

### Running as a container
This is an example of a containerized instance running in Docker. The following
example also leverages on local docker health check to restart the container
in case of internal failure. The check is executed every 10 minutes.
```
docker run -d -p 8080:8080 \
  --restart always \
  --name webcheckd \
  --health-cmd "curl --fail http://localhost:8080/healthz || exit 1"\
  --health-interval 5s \
  --health-retries 3 \
  quay.io/gbsalinetti/webcheckd \
  -url https://www.mysite.com \
  -interval 600 \
  -from admin@mysite.com \
  -password 'myrandompassword' \
  -to security-helpdesk@mysite.com \
  -to sre-helpdesk@mysite.com
```

### Running in Kubernetes
Additional files for K8S deployment are provided in the **k8s** directory.
First, create a secret holding the sender's account:
```
$ kubectl create secret generic webcheckd-secret \
  --from-literal=password=mypassw0rd \
  --from-literal=from=admin@mysite.com
```

Then, create a ConfigMap to provide URLs to check and notification recipients.
You can use the example ConfigMap provided in the **k8s** directory as a 
starting point.
```
$ kubectl apply -f k8s/cm.yaml
```

Then, create the deployment for the application using the deployment manifest
provided:
```
$ kubectl apply -f k8s/deployment.yaml
```

Finally, create a service to inspect the health check. Further features will
be enabled as REST API endpoints in the future.
```
$ kubectl apply -f k8s/svc.yaml
```

### Maintainer
Giovan Battista Salinetti <gbsalinetti@gmail.com>
