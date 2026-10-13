.PHONY: dev dev-web dev-server test test-go test-web lint lint-go lint-web build build-web build-go

dev:
	@echo "Run in two terminals: make dev-server / make dev-web"

dev-server:
	go run ./cmd/poker

dev-web:
	cd web && pnpm dev

test: test-go test-web

test-go:
	go test -race -cover ./...

test-web:
	cd web && pnpm test

lint: lint-go lint-web

lint-go:
	go vet ./...
	golangci-lint run

lint-web:
	cd web && pnpm lint && pnpm check

build: build-web build-go

build-web:
	cd web && pnpm install --frozen-lockfile && pnpm build
	find internal/web/dist -mindepth 1 ! -name .gitkeep -delete
	cp -R web/build/. internal/web/dist/

build-go:
	go build -o bin/poker ./cmd/poker
