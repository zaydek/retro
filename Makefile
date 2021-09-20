# This Makefile was authored on macOS and is unlikely to work on other operating
# systems. While the binaries this Makefile produces are highly portable, the
# development of these binaries is not generally portable.

# Get the working version from `version.txt`
VERSION = $(shell cat version.txt)

################################################################################

# Makes all binaries; `create-retro-app` and `retro`. Note that these binaries
# are moved to `/usr/local/bin` so that they may be tested locally. Aliasing
# these binaries is recommended for active development.
#
# ~/.bash_profile
#
# alias create-retro-app=/usr/local/bin/create-retro-app
# alias retro=/usr/local/bin/retro
#
all:
	make bin

# Makes `create-retro-app`
bin-create-retro-app:
	go build -o=create-retro-app main_create_retro_app.go && mv create-retro-app /usr/local/bin

# Makes `retro`
bin-retro:
	go build -o=retro main_retro.go && mv retro /usr/local/bin

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

# Builds Go binaries for platforms `darwin-64`, `linux-64`, `windows-64` and
# then transpiles `postinstall.ts` from TypeScript to JavaScript. Finally, a
# placeholder file is created so executables  `postinstall.js` and creates a
# placeholder file for `create-retro-app`.
build-create-retro-app:
	GOOS=darwin  GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/darwin-64 main_create_retro_app.go
	GOOS=linux   GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/linux-64 main_create_retro_app.go
	GOOS=windows GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/create-retro-app/bin/windows-64.exe main_create_retro_app.go

	./node_modules/.bin/esbuild \
		npm/create-retro-app/postinstall.ts \
			--format=cjs \
			--log-level=warning \
			--outfile=npm/create-retro-app/postinstall.esbuild.js \
			--target=es2018

	touch npm/create-retro-app/bin/create-retro-app

# Builds Go binaries for platforms `darwin-64`, `linux-64`, `windows-64` and
# then transpiles `postinstall.ts` from TypeScript to JavaScript. Finally, a
# placeholder file is created so executables  `postinstall.js` and creates a
# placeholder file for `retro`.
build-retro:
	./node_modules/.bin/esbuild \
		scripts/backend.ts \
			--format=cjs \
			--log-level=warning \
			--outfile=scripts/backend.esbuild.js \
			--target=es2018

	GOOS=darwin  GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/darwin-64 main_retro.go
	GOOS=linux   GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/linux-64 main_retro.go
	GOOS=windows GOARCH=amd64 go build "-ldflags=-s -w" -o=npm/retro/bin/windows-64.exe main_retro.go

	./node_modules/.bin/esbuild \
		npm/retro/postinstall.ts \
			--format=cjs \
			--log-level=warning \
			--outfile=npm/retro/postinstall.esbuild.js \
			--target=es2018

	touch npm/retro/bin/retro

	cp -r scripts npm/retro/bin && rm npm/retro/bin/scripts/backend.ts

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

clean:
	rm -rf scripts/backend.esbuild.js
	rm -rf npm/create-retro-app/postinstall.esbuild.js
	rm -rf npm/create-retro-app/bin
	rm -rf npm/retro/postinstall.esbuild.js
	rm -rf npm/retro/bin
