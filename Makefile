# ******************************************************************************
# Licensed Materials - Property of IBM
# IBM Cloud Schematics
# (C) Copyright IBM Corp. 2017 All Rights Reserved.
# US Government Users Restricted Rights - Use, duplication or
# disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
# ******************************************************************************

# Runs with out docker.
.DEFAULT_GOAL := run-mac

APPLICATION ?= $$(basename $(CURDIR))
DOCKER_REGISTRY ?= 'blueprint-docker-local.artifactory.swg-devops.com'
ORGANIZATION ?= 'blueprint'
GIT_VERSION := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
TRAVIS_BUILD_NUMBER ?= TestTravisBuildNo
MONGO ?= 'mongo'
MONGO_EXPRESS ?= 'mongo-express'
MONGO_PORT ?= 27017

export GONOSUMDB := 'github.ibm.com/*'
export GOPROXY := https://${ARTIFACTORY_USER_ID}:${ARTIFACTORY_PASSWORD}@na.artifactory.swg-devops.com/artifactory/api/go/blueprint-go-virtual

# Targets lint and check are needed?

.PHONY: vendor
vendor:
	go mod vendor

# Not tested
.PHONY: swagger
swagger:
	go generate ./server
	cd ./swagger; statik -src=./ui -f
	go mod vendor
	swagger generate spec -m -o ./swagger/ui/swagger.json
	swagger validate ./swagger/ui/swagger.json
	cd ./swagger; statik -src=./ui -f

.PHONY: environment
environment:
	$(eval API_HTTPADDR ?= "")
	$(eval API_HTTPPORT ?= 8080)
	$(eval MOUNT_DIR="/tmp")
	$(eval MONGO_DOCKER_HOST ?= "")
	$(eval MONGO_USER ?= root)
	$(eval MONGO_PASSWORD ?= example)

.PHONY: environment2
environment2:
	$(eval MONGO_HOST ?= $(shell docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${MONGO}))

# Refer for advanced usage
# https://www.bmc.com/blogs/mongodb-docker-container/
.PHONY: docker-run-mongo
docker-run-mongo: environment docker-stop-mongo
	docker run \
		-d \
		--name=${MONGO} \
		--hostname=${MONGO_DOCKER_HOST} \
		-p ${MONGO_PORT}:${MONGO_PORT} \
		-e MONGO_INITDB_ROOT_USERNAME=${MONGO_USER} \
		-e MONGO_INITDB_ROOT_PASSWORD=${MONGO_PASSWORD} \
		mongo:latest

.PHONY: docker-run-mongo-stack
docker-run-mongo-stack: environment docker-stop-mongo
	docker-compose \
	-f mongo-docker.yml up \
	--detach

# May not be needed - Erroring- Names cannot be given
.PHONY: docker-stop-mongo
docker-stop-mongo:
	docker stop ${MONGO} || true
	docker rm ${MONGO} || true

.PHONY: docker-attach-shell-mongo
docker-attach-shell-mongo:
	docker exec -it ${MONGO} /bin/bash

.PHONY: docker-build
docker-build: vendor
	SHA=$(shell git log -1 --pretty=format:"%H")
	go build -v -ldflags "-X main.commit=${GIT_VERSION} -X main.travisBuildNumber=${TRAVIS_BUILD_NUMBER} -X main.buildDate=${BUILD_DATE}"
	docker build \
	  --pull \
	  --no-cache \
	  -t ${APPLICATION} \
	  --build-arg gitSHA=$(GIT_VERSION) \
		--build-arg travisBuildNo=${TRAVIS_BUILD_NUMBER} \
		--build-arg buildDate=$(BUILD_DATE)\
	  .

.PHONY: docker-run
docker-run: environment docker-stop environment2
	docker run \
		-d \
		--name=${APPLICATION} \
		-p ${API_HTTPPORT}:${API_HTTPPORT} \
		-e API_MOUNT_DIR=${MOUNT_DIR} \
		-e API_HTTPPORT=${API_HTTPPORT} \
		-e API_HTTPADDR=${API_HTTPADDR} \
		-e API_MONGO_HOST=${MONGO_HOST} \
		-e API_MONGO_PORT=${MONGO_PORT} \
		-e API_MONGO_USERNAME=${MONGO_USER} \
		-e API_MONGO_PASSWORD=${MONGO_PASSWORD} \
		${APPLICATION}:latest

.PHONY: docker-stop
docker-stop:
	docker stop ${APPLICATION} || true
	docker rm ${APPLICATION} || true

.PHONY:sleep
sleep:
	sleep 5

.PHONY: run-local
run-local: environment environment2
	API_MOUNT_DIR=${MOUNT_DIR} \
	API_HTTPPORT=${API_HTTPPORT} \
	API_HTTPADDR=${API_HTTPADDR} \
	API_MONGO_HOST=${MONGO_HOST} \
	API_MONGO_PORT=${MONGO_PORT} \
	API_MONGO_USERNAME=${MONGO_USER} \
	API_MONGO_PASSWORD=${MONGO_PASSWORD} \
	go run .

.PHONY: run-mac
run-mac: environment environment2
	API_MOUNT_DIR=${MOUNT_DIR} \
	API_HTTPPORT=${API_HTTPPORT} \
	API_HTTPADDR=${API_HTTPADDR} \
	API_MONGO_HOST="localhost" \
	API_MONGO_PORT=${MONGO_PORT} \
	API_MONGO_USERNAME=${MONGO_USER} \
	API_MONGO_PASSWORD=${MONGO_PASSWORD} \
	go run .


#go install cmd/discovery/*.go
.PHONY: build-cli
build-cli: 
	cd cmd/discovery; go build

.PHONY: install-cli
install-cli:
	cd cmd/discovery; go install
	
