default: fmt test

test:
	go test *.go -v

fmt:
	go fmt *.go
