test:
	go test -v ./... --cover

deps:
	go get github.com/stretchr/testify 