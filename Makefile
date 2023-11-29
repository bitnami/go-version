.PHONY: all clean download get-build-deps vet lint test cover

include ./vars.mk

all:
	@$(MAKE) get-build-deps
	@$(MAKE) download
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) cover

clean:
	@rm -rf $(BUILD_DIR)

download:
	$(GO_MOD) download

get-build-deps:
	@echo "+ Downloading build dependencies"
	@go install honnef.co/go/tools/cmd/staticcheck@latest

vet:
	@echo "+ Vet"
	@go vet ./...

lint:
	@echo "+ Linting package"
	@staticcheck ./...
	$(call fmtcheck, .)

test:
	@echo "+ Testing package"
	$(GO_TEST) ./...

cover: test
	@echo "+ Tests Coverage"
	@mkdir -p $(BUILD_DIR)
	@touch $(BUILD_DIR)/cover.out
	@go test -coverprofile=$(BUILD_DIR)/cover.out ./... 
	@go tool cover -html=$(BUILD_DIR)/cover.out -o=$(BUILD_DIR)/coverage.html
