dependencies:
	go mod download

build-mocks:
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@go generate -v ./...

test:
	go test -tags testing ./...