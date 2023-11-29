SHELL := /bin/bash
GOPATH ?= $(shell go env GOPATH)
PATH := $(GOPATH)/bin:$(PATH)
BUILD_DIR := $(abspath ./out)

DEBUG ?= 0

ifeq ($(DEBUG),1)
GO_TEST := @go test -v
else
GO_TEST := @go test
endif

GO_MOD := @go mod
# Do not do goimport of the vendor dir
go_files=$$(find $(1) -type f -name '*.go' -not -path "./vendor/*")
fmtcheck = @if goimports -l $(go_files) | read var; then echo "goimports check failed for $(1):\n `goimports -d $(go_files)`"; exit 1; fi
