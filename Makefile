.PHONY: tailwind-watch
tailwind-watch:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

.PHONY: templ-build
templ-build: ## Build the templates
	templ generate

.PHONY: templ-watch
templ-watch: ## Watch for changes in the templates
	templ generate --watch

#.PHONY: dev-web
#dev-web:
#	 air

.PHONY: build-web
web: tailwind-build templ-build ## Build the web app binary
	go build -o ./bin/ilxweb ./cmd/web/main.go

.PHONY: build-indexer
indexer: templ-build vet ## Build the indexer binary
	go build -o ./bin/ilxindexer ./cmd/indexer/main.go

.PHONY: build
build: web indexer ## Build both binaries

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: clean
clean: ## Clean the binaries
	rm -rf ./bin

.PHONY: install-dev-dependencies
install-dev-dependencies: ## Install dev dependencies
	go install github.com/a-h/templ/cmd/templ@latest

.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
