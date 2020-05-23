# Install all the tools for development
.PHONY: init
init: lint-prepare goimports-prepare

.PHONY: lint-prepare
lint-prepare:
	@echo "Preparing Linter"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

.PHONY: goimports-prepare
goimports-prepare:
	@echo "Preparing goimports"
	go get golang.org/x/tools/cmd/goimports
	go mod tidy
