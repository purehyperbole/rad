test:
	go test -v --race ./...

bench:
	gotest -v -bench=. -benchmem -benchtime=1000000x

deps:
	go get github.com/stretchr/testify
	go get github.com/google/uuid
