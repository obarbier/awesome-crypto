.PHONY: go-generate test

GO_CMD?=go

bootstrap:
	$(GO_CMD) generate -tags tools tools/tools.go

go-generate:
	$(GO_CMD) generate ./domain/ports.go

test:
	$(GO_CMD) test