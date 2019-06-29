test:
	go test -v --race ./...

deps:
	go get github.com/stretchr/testify
	go get github.com/google/uuid