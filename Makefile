VERSION = $(shell cat version.txt)


################################################################################

all:
	make bin

bin-create-retro-app:
	go build -o=create-retro-app main_create_retro_app.go && mv create-retro-app /usr/local/bin

bin-retro:
	go build -o=retro main_retro.go && mv retro /usr/local/bin

bin:
	make -j2 \
		bin-create-retro-app \
		bin-retro

################################################################################

test-create-retro-app:
	go test ./cmd/create_retro_app/...

test-retro:
	go test ./cmd/retro/...

test-pkg:
	go test ./pkg/...

test:
	go test ./...

################################################################################

build-create-retro-app:
	GOOS=darwin  GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/darwin-64 main_create_retro_app.go
	GOOS=linux   GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/linux-64 main_create_retro_app.go
	GOOS=windows GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/windows-64.exe main_create_retro_app.go

	./node_modules/.bin/esbuild \
		npm/create-retro-app/postinstall.ts \
			--format=cjs \
			--log-level=warning \
			--outfile=npm/create-retro-app/postinstall.esbuild.js
	touch npm/create-retro-app/bin/create-retro-app

build-retro:
	./node_modules/.bin/esbuild \
		scripts/backend.ts \
			--format=cjs \
			--log-level=warning \
			--outfile=scripts/backend.esbuild.js

	GOOS=darwin  GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/darwin-64 main_retro.go
	GOOS=linux   GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/linux-64 main_retro.go
	GOOS=windows GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/windows-64.exe main_retro.go

	./node_modules/.bin/esbuild \
		npm/retro/postinstall.ts \
			--format=cjs \
			--log-level=warning \
			--outfile=npm/retro/postinstall.esbuild.js
	touch npm/retro/bin/retro

	cp -r scripts npm/retro/bin && rm npm/retro/bin/scripts/backend.ts

build:
	make -j2 \
		build-create-retro-app \
		build-retro

################################################################################

version:
	cd npm/create-retro-app && npm version "$(VERSION)" --allow-same-version
	cd npm/retro && npm version "$(VERSION)" --allow-same-version

################################################################################

release-dry-run:
	cd npm/create-retro-app && npm publish --dry-run
	cd npm/retro && npm publish --dry-run

release:
	cd npm/create-retro-app && npm publish
	cd npm/retro && npm publish

################################################################################

clean:
	rm -rf scripts/backend.esbuild.js
	rm -rf npm/create-retro-app/postinstall.esbuild.js
	rm -rf npm/create-retro-app/bin
	rm -rf npm/retro/postinstall.esbuild.js
	rm -rf npm/retro/bin
