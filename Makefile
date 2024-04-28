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

.PHONY: web
web: tailwind-build templ-build ## Build the web app binary
	go build -o ./bin/ilx-web ./cmd/web/main.go

.PHONY: indexer
indexer: ## Build the indexer binary
	go build -o ./bin/ilx-indexer ./cmd/indexer/main.go

.PHONY: docker-web
docker-web:  ## Build the web app docker image
	docker build -t tylersmith/ilx-web -f docker/web.dockerfile .

.PHONY: docker-indexer
docker-indexer:  ## Build the indexer docker image
	docker build -t tylersmith/ilx-indexer -f docker/indexer.dockerfile .

.PHONY: clean
clean: ## Clean the binaries
	rm -rf ./bin ./tmp

.PHONY: build-deps
build-deps: ## Install dev dependencies
	go install github.com/a-h/templ/cmd/templ@latest

.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
