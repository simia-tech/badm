
GO ?= go

all:
	$(GO) build -o badm cmd/badm/main.go
