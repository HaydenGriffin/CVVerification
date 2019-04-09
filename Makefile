.PHONY: all dev clean build env-up env-down run

all: clean build env-up run init

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@cd chaincode && dep ensure
	@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fabric-network/fixtures && docker-compose up -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fabric-network/fixtures && docker-compose down
	@echo "Environment down"

##### RUN
init:
	@echo "Start app and init ..."
	@cd app && ./app -install -register

run:
	@echo "Start app ..."
	@./cvverification

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/cvverification-* app/app chaincode/chaincode
	@docker rm -f -v `docker ps -a --no-trunc | grep "cvverification" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "cvverification" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
