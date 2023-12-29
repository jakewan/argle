.DEFAULT_GOAL := all

.PHONY: all
all: go-fmt go-test go-vet go-staticcheck

.PHONY: go-fmt
go-fmt:
	gofmt -d -s -w .

.PHONY: go-staticcheck
go-staticcheck:
	staticcheck -go 1.21 ./...

.PHONY: go-test
go-test:
	go test -v ./...

.PHONY: go-vet
go-vet:
	go vet ./...
