.PHONY: go-generate test

bootstrap:
	go generate -tags tools tools/tools.go

go-generate:
	go generate ./domain/ports.go
test:
	go test