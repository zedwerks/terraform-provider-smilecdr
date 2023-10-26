TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=zedwerks
NAME=smilecdr
OUTPUT_DIR=./bin
BINARY=${OUTPUT_DIR}/terraform-provider-${NAME}_v${VERSION}
VERSION=$$(git describe --tags)
OS_ARCH?=darwin_arm64
DIST_DIR=./dist
BINARY=terraform-provider-${NAME}_${VERSION}
VERSION=$$(git describe --tags)
OS_ARCH?=darwin_arm64

default: docs build

build: show-version
	mkdir -p ${OUTPUT_DIR}
	@echo "Building ${OUTPUT_DIR}/${BINARY}"
	go build -o ${OUTPUT_DIR}/${BINARY}

docs:
	@echo "Generating docs"
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

show-version:
	git describe --tags

release:
	goreleaser release --clean --snapshot --skip=publish  --skip=sign

binaries: build
	GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/${BINARY}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ${OUTPUT_DIR}/${BINARY}_darwin_arm64
	GOOS=freebsd GOARCH=amd64 go build -o ${OUTPUT_DIR}/${BINARY}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ${OUTPUT_DIR}/${BINARY}_freebsd_arm
	GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/${BINARY}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ${OUTPUT_DIR}/${BINARY}_linux_arm
	GOOS=openbsd GOARCH=amd64 go build -o ${OUTPUT_DIR}/${BINARY}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ${OUTPUT_DIR}/${BINARY}_solaris_amd64
	GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/${BINARY}_windows_amd64

install: build 
	mkdir -p ~/.terraform.d/plugins/local.providers/${HOSTNAME}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/local.providers/${HOSTNAME}/${NAME}/${VERSION}/${OS_ARCH}


test: 
	go test -i $(TEST) || exit 1                                                   
	@echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc:
	@echo "Running acceptance tests"
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}
