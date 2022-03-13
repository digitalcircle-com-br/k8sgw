GIT_COMMIT := $(shell git rev-list -1 HEAD)
DT := $(shell date +%Y.%m.%d.%H%M%S)
ME := $(shell whoami)
HOST := $(shell hostname)
PRODUCT := k8sgw
AWS_REGION := sa-east-1
AWS_ECR := 294726603466.dkr.ecr.sa-east-1.amazonaws.com


ver:
	echo package main > ver.go
	echo const VER=\"$(DT)\" >> ver.go

run:
	CGO_ENABLED=0 go run -ldflags "-X github.com/digitalcircle-com-br/buildinfo.Ver=$(GIT_COMMIT) -X github.com/digitalcircle-com-br/buildinfo.BuildDate=$(DT) -X github.com/digitalcircle-com-br/buildinfo.BuildUser=$(ME) -X github.com/digitalcircle-com-br/buildinfo.BuildHost=$(HOST) -X github.com/digitalcircle-com-br/buildinfo.Product=$(PRODUCT)" ./main.go

docker_push: docker
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ECR) && \
	docker tag $(PRODUCT):latest $(AWS_ECR)/$(PRODUCT):latest && \
	docker push $(AWS_ECR)/$(PRODUCT):latest

docker: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deploy/k8sgw/k8sgw -ldflags "-s -w -X github.com/digitalcircle-com-br/buildinfo.Ver=$(GIT_COMMIT) -X github.com/digitalcircle-com-br/buildinfo.BuildDate=$(DT) -X github.com/digitalcircle-com-br/buildinfo.BuildUser=$(ME) -X github.com/digitalcircle-com-br/buildinfo.BuildHost=$(HOST) -X github.com/digitalcircle-com-br/buildinfo.Product=$(PRODUCT)" ./main.go
	cd deploy/k8sgw && \
	docker build -t $(PRODUCT) .
	#docker tag $(PRODUCT) 10.100.0.10:5050/$(PRODUCT)
	#docker push  10.100.0.10:5050/$(PRODUCT)

docker_run:
	docker run --rm -it -p 8080:8080 $(AWS_ECR)/$(PRODUCT):latest

sample_build:
	CGO_ENABLED=0 go run -ldflags "-X github.com/digitalcircle-com-br/buildinfo.Ver=$(GIT_COMMIT) -X github.com/digitalcircle-com-br/buildinfo.BuildDate=$(DT) -X github.com/digitalcircle-com-br/buildinfo.BuildUser=$(ME) -X github.com/digitalcircle-com-br/buildinfo.BuildHost=$(HOST) -X github.com/digitalcircle-com-br/buildinfo.Product=sampleapp" ./deploy/sampleapp/main.go
	cd deploy/sampleapp && docker build -t sampleapp .
	docker tag sampleapp 294726603466.dkr.ecr.sa-east-1.amazonaws.com/sampleapp
	docker push 294726603466.dkr.ecr.sa-east-1.amazonaws.com/sampleapp
