# This Makefile was authored on macOS and is unlikely to work on other operating
# systems. While the binaries this Makefile produces are highly portable, the
# development of these binaries is not generally portable.

# Get the working version from `version.txt`
VERSION = $(shell cat version.txt)

################################################################################

# Bundles the Node.js backend code
bundle:
	npx esbuild scripts/backend/backend.ts \
		--bundle \
		--external:esbuild --external:react --external:react-dom --external:react-dom/server \
		--log-level=warning \
		--outfile=scripts/backend.esbuild.js \
		--platform=node \
		--sourcemap

################################################################################

# Makes all binaries; `create-retro-app` and `retro`. Note that these binaries
# are moved to `~/github/bin` so that they may be tested locally. Aliasing these
# binaries is recommended for active development.
#
# ~/.bash_profile
#
# alias create-retro-app=~/github/bin/create-retro-app
# alias retro=~/github/bin/retro
#
all:
	make bin

# Makes `create-retro-app`
bin-create-retro-app:
	go build -o=create-retro-app main_create_retro_app.go && mv create-retro-app ~/github/bin

# Makes `retro`
bin-retro:
	make bundle
	go build -o=retro main_retro.go && mv retro ~/github/bin

# Makes all binaries in parallel
bin:
	make -j2 \
		bin-create-retro-app \
		bin-retro

################################################################################

# Run all Go tests for `create-retro-app`
test-create-retro-app:
	go test ./cmd/create_retro_app/...

# Run all Go tests for `retro`
test-retro:
	go test ./cmd/retro/...

# Run all Go tests for local dependencies
test-pkg:
	go test ./pkg/...

# Run all Go tests (not in parallel)
test:
	make test-create-retro-app
	make test-retro
	make test-pkg

################################################################################

# Builds Go binaries and creates a placeholder executable for the post-
# installation script
build-create-retro-app:
	rm -rf npm/create_retro_app/bin
	mkdir -p npm/create_retro_app/bin

	GOOS=darwin  GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/darwin-64      main_create_retro_app.go
	GOOS=linux   GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/linux-64       main_create_retro_app.go
	GOOS=windows GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/windows-64.exe main_create_retro_app.go

	touch npm/create-retro-app/bin/create-retro-app

# Builds Go binaries and creates a placeholder executable for the post-
# installation script
build-retro:
	rm -rf npm/retro/bin
	mkdir -p npm/retro/bin/scripts

	make bundle

	GOOS=darwin  GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/darwin-64      main_retro.go
	GOOS=linux   GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/linux-64       main_retro.go
	GOOS=windows GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/windows-64.exe main_retro.go

	cp \
		scripts/backend.esbuild.js \
		scripts/backend.esbuild.js.map \
		scripts/require.js \
		scripts/vendor.js \
		npm/retro/bin/scripts
	touch npm/retro/bin/retro

# Makes all builds in parallel
build:
	make -j2 \
		build-create-retro-app \
		build-retro

################################################################################

# Versions `create-retro-app` and `retro`
version:
	cd npm/create-retro-app && npm version "$(VERSION)" --allow-same-version
	cd npm/retro && npm version "$(VERSION)" --allow-same-version

################################################################################

# Releases (dry-run) `create-retro-app` and `retro`
release-dry-run:
	cd npm/create-retro-app && npm publish --dry-run
	cd npm/retro && npm publish --dry-run

# Releases `create-retro-app` and `retro`
release:
	cd npm/create-retro-app && npm publish
	cd npm/retro && npm publish

################################################################################

# Publishes (dry-run) `create-retro-app` and `retro`
publish-dry-run:
	make build && \
		make version && \
		make release-dry-run

# Publishes `create-retro-app` and `retro`
publish:
	make build && \
		make version && \
		make release
