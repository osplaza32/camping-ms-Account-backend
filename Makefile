.DEFAULT_GOAL := help

# Show this help.
help:           ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

run: ## run Grpc server and proxy whithout docker
	go run cmd/main.go

up: ## up docker compose
	docker-compose up -d

down: ## down docker compose
	docker-compose down
vendor: ## dowload all dependecies
	go mod vendor
	cd proto && glide install
# general
vendor-proto: ## dowload all dependecies
	cd proto && glide install

# coverage

CURRENT_TAG := $(shell git rev-parse --short HEAD)
WORKDIR = $(PWD)
DOCKER_FILE = $(WORKDIR)/Dockerfile
current_dir = $(notdir $(shell pwd))
CUR_DIR = $(shell echo $(current_dir) | tr A-Z a-z)
DOCKER_REPO = gcr.io/$(shell gcloud config list --format 'value(core.project)')/osplaza32/$(CUR_DIR)
COVERAGE_REPORT = coverage.txt
COVERAGE_PROFILE = profile.out
COVERAGE_MODE = atomic


MY_VAR := $(shell find $(WORKDIR) -name "*.go" -not -path "*/vendor/*" -not -path "*/gen/*" | grep -o '.*/' |uniq  | xargs)

coverage: ## Coverage test
	@cd $(WORKDIR); \
	go test -v  -race ${MY_VAR} -coverprofile=$(COVERAGE_PROFILE) -covermode=$(COVERAGE_MODE);
compile: ## compile protobuffer
	@cd $(WORKDIR)/proto/src; \
	prototool cache delete;
	@cd $(WORKDIR)/proto/src; \
	prototool generate;
all-docker: ## init all and play witch docker composer
	make vendor-proto && make compile && make vendor && make up ;
all: ## init all and play witch docker composer
	make vendor-proto && make compile && make vendor && make run ;
stop: ## delete all docker composer
	make down;
all-dep: ## init all and play witch docker composer
	make vendor-proto && make compile && make vendor;

image-deployment-first:
	docker build . -f $(DOCKER_FILE) -t $(DOCKER_REPO):$(CURRENT_TAG)
	docker tag $(DOCKER_REPO):$(CURRENT_TAG) $(DOCKER_REPO):latest
	docker push $(DOCKER_REPO):$(CURRENT_TAG)
	kubectl create deployment automate-grcp-server --image=$(DOCKER_REPO):$(CURRENT_TAG)
	kubectl expose deployment automate-grcp-server --type LoadBalancer --name=http --port 80 --target-port 8080
	kubectl expose deployment automate-grcp-server --type LoadBalancer --name=grcp --port 50050 --target-port 50051
	kubectl get pods

image-deployment:
	docker build . -f $(DOCKER_FILE) -t $(DOCKER_REPO):$(CURRENT_TAG)
	docker tag $(DOCKER_REPO):$(CURRENT_TAG) $(DOCKER_REPO):latest
	docker push $(DOCKER_REPO):$(CURRENT_TAG)
	kubectl set image deployment/automate-grcp-server --image=image:$(DOCKER_REPO):$(CURRENT_TAG)
	kubectl expose deployment automate-grcp-server --type LoadBalancer --name=http --port 80 --target-port 8080 && kubectl expose deployment automate-grcp-server --type LoadBalancer --name=grcp --port 50050 --target-port 50051
register:
	docker build . -f $(DOCKER_FILE) -t $(DOCKER_REPO):$(CURRENT_TAG)
	docker tag $(DOCKER_REPO):$(CURRENT_TAG) $(DOCKER_REPO):latest
	docker push $(DOCKER_REPO):$(CURRENT_TAG)
print:
	$(CUR_DIR)
