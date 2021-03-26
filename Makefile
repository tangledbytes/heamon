GIT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_COMMITSHA=$(shell git rev-list -1 HEAD)

run-local:
	PORT=5000 go run ./cmd/rest/main.go;

run-ui:
	cd ui && npm start

build: build-ui build-server

build-server:
	go build -ldflags="-w -s -X main.version=${GIT_VERSION} -X main.commit=${GIT_COMMITSHA}" -a -o heamon ./cmd/rest/main.go

build-ui:
	cd ui && npm run build && rm -rf ../pkg/entrypoint/rest/ui && mv build ../pkg/entrypoint/rest/ui

setup-ui:
	cd ui && npm i