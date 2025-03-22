TEST?=$$(go list ./... | grep -v 'vendor' | grep -v 'env' | grep -v 'utils')
HOSTNAME=t1
NAMESPACE=t1-cloud
NAME=t1
VERSION=1.6.1
BINARY=server
OS_ARCH=darwin_amd64
TERRAFORM_VERSION=1.9.2
TF_LOG_FILE=terraform_provider_t1_runtime.log

default: build

build: build_server build_agent

build_server:
	go build -o ${BINARY} cmd/server/main.go

build_agent:
	go build -o ${BINARY} cmd/agent/main.go

test_server:
	go test ./... -timeout 60m --tags=server -v

test_agent:
	go test ./... -timeout 60m --tags=agent -v
