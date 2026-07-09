.PHONY: check build frontend backend clean

check:
	cd frontend && npm ci && npm run check
	cd backend && go build ./cmd/server

build: backend frontend

backend:
	cd backend && go build -trimpath -ldflags="-s -w" -o ../dist/caro-server ./cmd/server

frontend:
	cd frontend && npm ci && npm run build
	mkdir -p dist/frontend
	cp -R frontend/build/. dist/frontend/

clean:
	rm -rf dist frontend/build frontend/.svelte-kit
