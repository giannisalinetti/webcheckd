APP_NAME=webcheckd
REGISTRY=quay.io
NAMESPACE=gbsalinetti

build:
	GO111MODULE=on go build .
image:
	docker build -t $(APP_NAME) .

tag:
	docker tag $(APP_NAME) $(REGISTRY)/$(NAMESPACE)/$(APP_NAME):latest

push:
	docker push $(REGISTRY)/$(NAMESPACE)/$(APP_NAME):latest

