#
# @file
# Footle's Makefile.
#
# Known to work with GNU make in GNU/Linux.
#
# Puts all build output into the "build/" directory.  Successful builds should
# produce *at least* an executable at "build/bin/footle".
#
# Required tools:
# - Go >= 1.8: https://golang.org/doc/install
# - Go dep: https://github.com/golang/dep#installation
# - go-bindata: go get -u https://github.com/go-bindata/go-bindata/...
# - Node.js and npm: https://github.com/creationix/nvm
#
# The dep and go-bindata executables *must* be in $PATH
#
# Useful make commands:
# - make: Compiles Footle for the current platform.
# - make dist: Cross compiles for FreeBSD, Linux, Mac, and Windows.
# - make clean: Removes build for the current platform.
# - make distclean: Removes all cross compiled builds.
# - make realclean: Removes all types of builds.

BUILD_ROOT = ./build

# Build directory for the Footle binary for the current platform.
BIN_BUILD_DIR_PATH = ${BUILD_ROOT}/bin

# Cross platform builds go here.
DIST_DIR = ${BUILD_ROOT}/dist

# Build directories for all platforms.
LINUX_32_BUILD_DIR = footle-linux-32
LINUX_64_BUILD_DIR = footle-linux-64
FREEBSD_64_BUILD_DIR = footle-freebsd-64
MACOS_64_BUILD_DIR   = footle-mac-64
WIN_32_BUILD_DIR = footle-win-32
WIN_64_BUILD_DIR = footle-win-64
LINUX_32_BUILD_DIR_PATH = ${DIST_DIR}/${LINUX_32_BUILD_DIR}
LINUX_64_BUILD_DIR_PATH = ${DIST_DIR}/${LINUX_64_BUILD_DIR}
FREEBSD_64_BUILD_DIR_PATH = ${DIST_DIR}/${FREEBSD_64_BUILD_DIR}
MACOS_64_BUILD_DIR_PATH   = ${DIST_DIR}/${MACOS_64_BUILD_DIR}
WIN_32_BUILD_DIR_PATH = ${DIST_DIR}/${WIN_32_BUILD_DIR}
WIN_64_BUILD_DIR_PATH = ${DIST_DIR}/${WIN_64_BUILD_DIR}

# Location of Go code.
SERVER_SRC_DIR_PATH = ./src/server
GO_SRC_FILES = $(shell find ${SERVER_SRC_DIR_PATH} -name *.go -type f)
GO_UI_BUNDLE_DIRECTIVE_FILE = ${SERVER_SRC_DIR_PATH}/http/http.go
GO_UI_BUNDLE_SRC_FILE = ${SERVER_SRC_DIR_PATH}/http/uibundle/ui_bundle.go
OUR_GO_PATH = $(shell pwd)

# Location of HTML/CSS/Javascript-based UI code.
UI_SRC_DIR_PATH = ./src/ui
UI_BUILD_DIR_PATH  = ${BUILD_ROOT}/ui
UI_SRC_FILES = $(shell find ${UI_SRC_DIR_PATH} -path ${UI_SRC_DIR_PATH}/node_modules -prune -o -type f)

UI_HTML_SRC_PATH = ${UI_SRC_DIR_PATH}/index.html
UI_HTML_BUILD_PATH = ${UI_BUILD_DIR_PATH}/index.html

UI_SASS_DIR = ${UI_SRC_DIR_PATH}/style/sass
UI_SASS_BUILD_DIR = ${UI_BUILD_DIR_PATH}/style/sass
UI_CSS_DIR = ${UI_BUILD_DIR_PATH}/style/css
UI_SASS_FILES = $(shell find ${UI_SASS_DIR} -type f)

UI_SCRIPT_SRC_DIR_PATH = ${UI_SRC_DIR_PATH}/scripts
UI_SCRIPT_BUILD_DIR_PATH = ${UI_BUILD_DIR_PATH}/scripts
UI_SCRIPT_SRC_FILES = $(shell find ${UI_SCRIPT_SRC_DIR_PATH} -type f)

UI_FONT_ORIG_DIR = ${UI_BUILD_DIR_PATH}/node_modules/uikit/dist/fonts
UI_FONT_REQUIRED_DIR = ${UI_BUILD_DIR_PATH}/style/fonts


# Non-file targets.
.PHONY: all server embedded-ui godeps ui ui-dependencies markup style script font test test-execution dist cross-compile linux32 linux64 freebsd64 macos64 win32 win64 doc-copy doc-copy-linux32 doc-copy-linux64 doc-copy-freebsd64 doc-copy-macos64 doc-copy-win32 doc-copy-win64 tarball tarball-linux32 tarball-linux64 tarball-freebsd64 tarball-macos64 zipball-win32 zipball-win64 realclean distclean clean


# Compile Footle server's Go code and Footle UI's CSS/HTML/Javascript code.
all: ui server


# Compile the Footle binary.
server: godeps embedded-ui ${BIN_BUILD_DIR_PATH}/footle

# *Temporarily* set $GOPATH to Footle's top level dir so that Footle's Go
# packages (e.g. server/config) can be located during compilation.  This
# frees Footle from the Go workspace allowing it to reside anywhere
# in the file system.
${BIN_BUILD_DIR_PATH}/footle: ${GO_UI_BUNDLE_SRC_FILE} ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} go build -o ${BIN_BUILD_DIR_PATH}/footle ${SERVER_SRC_DIR_PATH}

# Grab Go package dependencies.
#
# We are using "dep" as the dependency manager.  "dep" demands that Footle be
# inside $GOPATH.  So we temporarily point $GOPATH to Footle's top level dir.
# This has the side effect that "dep" drops a "pkg/" directory inside Footle.
# Once "dep" is finished, we remove "pkg/" for the sake of cleanliness.
godeps: ${SERVER_SRC_DIR_PATH}/vendor
${SERVER_SRC_DIR_PATH}/vendor: ${SERVER_SRC_DIR_PATH}/Gopkg.toml
	cd ${SERVER_SRC_DIR_PATH} ; \
	GOPATH=${OUR_GO_PATH} dep ensure
	rm -rf ${OUR_GO_PATH}/pkg

# Prepare a Go source file containing all the UI files.  This is how we embed
# the UI files into the footle binary.
embedded-ui: ui ${GO_UI_BUNDLE_SRC_FILE}
# This generate directive relies on ./build/ui, but we are not using it here as
# a dependency.  The assumption is that whenever ${GO_SRC_FILES}
# changes, ./build/ui will also be updated by the "ui" target.
${GO_UI_BUNDLE_SRC_FILE}: ${UI_SRC_FILES}
	go generate ${GO_UI_BUNDLE_DIRECTIVE_FILE}


### Compile and prepare the UI #################################################
# Prepare HTML-based UI.
ui: ui-dependencies markup style script font

# Node.js packages.
ui-dependencies: ${UI_SRC_DIR_PATH}/node_modules ${UI_BUILD_DIR_PATH}/node_modules

${UI_SRC_DIR_PATH}/node_modules: ${UI_SRC_DIR_PATH}/package-lock.json
	cd ${UI_SRC_DIR_PATH}; \
	npm install
	touch $@

${UI_BUILD_DIR_PATH}/node_modules: ${UI_BUILD_DIR_PATH}/package.json ${UI_BUILD_DIR_PATH}/package-lock.json
	cd ${UI_BUILD_DIR_PATH}; \
	npm install --production
	touch $@

${UI_BUILD_DIR_PATH}/package.json: ${UI_SRC_DIR_PATH}/package.json
	mkdir -p ${UI_BUILD_DIR_PATH}
	cp $< $@

${UI_BUILD_DIR_PATH}/package-lock.json: ${UI_SRC_DIR_PATH}/package-lock.json
	mkdir -p ${UI_BUILD_DIR_PATH}
	cp $< $@

# Prepare CSS whenever *any* Sass file changes.
style: ui-dependencies ${UI_CSS_DIR}/ui.css
${UI_CSS_DIR}/ui.css: ${UI_SASS_FILES}
	# For the CSS Sourcemap to be of any use, our Sass files also have to be
	# served by the web server.
	mkdir -p ${UI_SASS_BUILD_DIR}
	cp -r ${UI_SASS_DIR}/* ${UI_SASS_BUILD_DIR}
	mkdir -p ${UI_CSS_DIR}
	cd ${UI_SRC_DIR_PATH}; \
	npx --no-install node-sass --indented-syntax --source-map true ../../${UI_SASS_BUILD_DIR}/ui.sass ../../$@

# Copy index.html
markup: ${UI_HTML_BUILD_PATH}
${UI_HTML_BUILD_PATH}: ${UI_HTML_SRC_PATH}
	cp $< $@

# Copy Javascript whenever *any* script file changes.
script: ${UI_SCRIPT_BUILD_DIR_PATH} 
${UI_SCRIPT_BUILD_DIR_PATH}: $(wildcard ${UI_SCRIPT_SRC_DIR_PATH} ${UI_SCRIPT_SRC_DIR_PATH}/*)
	cp -r ${UI_SCRIPT_SRC_DIR_PATH}/. $@
	touch $@

# Relocate uikit's "fonts" directory.  Otherwise it won't be found by the
# web browser.
font: ui-dependencies ${UI_FONT_REQUIRED_DIR}
${UI_FONT_REQUIRED_DIR}: ${UI_FONT_ORIG_DIR}
	cp -r $< $@


### Run tests ##################################################################
#
# Run all unit tests.
#
# At the moment, we have tests for the Go code only.
#
# We are explicitely mentioning the Go packages because we do *not* want to
# test the vendor packages.  Go 1.9 ignores the vendor packages.  So this is
# a temporary arrangement while we support Go 1.8.
test: embedded-ui test-execution
test-execution:
	GOPATH=${OUR_GO_PATH} go test ${SERVER_SRC_DIR_PATH}/core/... ${SERVER_SRC_DIR_PATH}/dbgp/... ${SERVER_SRC_DIR_PATH}/http/...


### Cross compile and prepare tarballs ########################################
#
# Prepare distributions for various platforms.
#
# Preparation formula: Cross compile Go code, copy UI code, create tar/zip
# archive.
#
dist: cross-compile doc-copy tarball

cross-compile: embedded-ui linux32 linux64 freebsd64 macos64 win32 win64

linux32: godeps ${LINUX_32_BUILD_DIR_PATH}/footle
${LINUX_32_BUILD_DIR_PATH}/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=linux GOARCH=386 go build -o $@ ${SERVER_SRC_DIR_PATH}

linux64: godeps ${LINUX_64_BUILD_DIR_PATH}/footle
${LINUX_64_BUILD_DIR_PATH}/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=linux GOARCH=amd64 go build -o $@ ${SERVER_SRC_DIR_PATH}

freebsd64: godeps ${FREEBSD_64_BUILD_DIR_PATH}/footle
${FREEBSD_64_BUILD_DIR_PATH}/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=freebsd GOARCH=amd64 go build -o $@ ${SERVER_SRC_DIR_PATH}

macos64: godeps ${MACOS_64_BUILD_DIR_PATH}/footle
${MACOS_64_BUILD_DIR_PATH}/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=darwin GOARCH=amd64 go build -o $@  ${SERVER_SRC_DIR_PATH}

win32: godeps ${WIN_32_BUILD_DIR_PATH}/footle.exe
${WIN_32_BUILD_DIR_PATH}/footle.exe: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=windows GOARCH=386 go build -o $@ ${SERVER_SRC_DIR_PATH}

win64: godeps ${WIN_64_BUILD_DIR_PATH}/footle.exe
${WIN_64_BUILD_DIR_PATH}/footle.exe: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=windows GOARCH=amd64 go build -o $@ ${SERVER_SRC_DIR_PATH}


# Copy the README and LICENSE files inside every distribution.
doc-copy: cross-compile doc-copy-linux32 doc-copy-linux64 doc-copy-freebsd64 doc-copy-macos64 doc-copy-win32 doc-copy-win64

doc-copy-linux32: ${LINUX_32_BUILD_DIR_PATH}/README.md ${LINUX_32_BUILD_DIR_PATH}/LICENSE
${LINUX_32_BUILD_DIR_PATH}/README.md: ${OUR_GO_PATH}/README.md
	cp ${OUR_GO_PATH}/README.md ${LINUX_32_BUILD_DIR_PATH}/README.md
${LINUX_32_BUILD_DIR_PATH}/LICENSE: ${OUR_GO_PATH}/LICENSE
	cp ${OUR_GO_PATH}/LICENSE ${LINUX_32_BUILD_DIR_PATH}/LICENSE

doc-copy-linux64: ${LINUX_64_BUILD_DIR_PATH}/README.md ${LINUX_64_BUILD_DIR_PATH}/LICENSE
${LINUX_64_BUILD_DIR_PATH}/README.md: ${OUR_GO_PATH}/README.md
	cp ${OUR_GO_PATH}/README.md ${LINUX_64_BUILD_DIR_PATH}/README.md
${LINUX_64_BUILD_DIR_PATH}/LICENSE: ${OUR_GO_PATH}/LICENSE
	cp ${OUR_GO_PATH}/LICENSE ${LINUX_64_BUILD_DIR_PATH}/LICENSE

doc-copy-freebsd64: ${FREEBSD_64_BUILD_DIR_PATH}/README.md ${FREEBSD_64_BUILD_DIR_PATH}/LICENSE
${FREEBSD_64_BUILD_DIR_PATH}/README.md: ${OUR_GO_PATH}/README.md
	cp ${OUR_GO_PATH}/README.md ${FREEBSD_64_BUILD_DIR_PATH}/README.md
${FREEBSD_64_BUILD_DIR_PATH}/LICENSE: ${OUR_GO_PATH}/LICENSE
	cp ${OUR_GO_PATH}/LICENSE ${FREEBSD_64_BUILD_DIR_PATH}/LICENSE

doc-copy-macos64: ${MACOS_64_BUILD_DIR_PATH}/README.md ${MACOS_64_BUILD_DIR_PATH}/LICENSE
${MACOS_64_BUILD_DIR_PATH}/README.md: ${OUR_GO_PATH}/README.md
	cp ${OUR_GO_PATH}/README.md ${MACOS_64_BUILD_DIR_PATH}/README.md
${MACOS_64_BUILD_DIR_PATH}/LICENSE: ${OUR_GO_PATH}/LICENSE
	cp ${OUR_GO_PATH}/LICENSE ${MACOS_64_BUILD_DIR_PATH}/LICENSE

doc-copy-win32: ${WIN_32_BUILD_DIR_PATH}/README.md ${WIN_32_BUILD_DIR_PATH}/LICENSE
${WIN_32_BUILD_DIR_PATH}/README.md: ${OUR_GO_PATH}/README.md
	cp ${OUR_GO_PATH}/README.md ${WIN_32_BUILD_DIR_PATH}/README.md
${WIN_32_BUILD_DIR_PATH}/LICENSE: ${OUR_GO_PATH}/LICENSE
	cp ${OUR_GO_PATH}/LICENSE ${WIN_32_BUILD_DIR_PATH}/LICENSE

doc-copy-win64: ${WIN_64_BUILD_DIR_PATH}/README.md ${WIN_64_BUILD_DIR_PATH}/LICENSE
${WIN_64_BUILD_DIR_PATH}/README.md: ${OUR_GO_PATH}/README.md
	cp ${OUR_GO_PATH}/README.md ${WIN_64_BUILD_DIR_PATH}/README.md
${WIN_64_BUILD_DIR_PATH}/LICENSE: ${OUR_GO_PATH}/LICENSE
	cp ${OUR_GO_PATH}/LICENSE ${WIN_64_BUILD_DIR_PATH}/LICENSE


tarball: cross-compile doc-copy tarball-linux32 tarball-linux64 tarball-freebsd64 tarball-macos64 zipball-win32 zipball-win64

tarball-linux32: ${DIST_DIR}/${LINUX_32_BUILD_DIR}.tar.gz
${DIST_DIR}/${LINUX_32_BUILD_DIR}.tar.gz: $(wildcard ${LINUX_32_BUILD_DIR_PATH} ${LINUX_32_BUILD_DIR_PATH}/*)
	cd ${DIST_DIR} ; tar cfz ${LINUX_32_BUILD_DIR}.tar.gz ${LINUX_32_BUILD_DIR}

tarball-linux64: ${DIST_DIR}/${LINUX_64_BUILD_DIR}.tar.gz
${DIST_DIR}/${LINUX_64_BUILD_DIR}.tar.gz: $(wildcard ${LINUX_64_BUILD_DIR_PATH} ${LINUX_64_BUILD_DIR_PATH}/*)
	cd ${DIST_DIR} ; tar cfz ${LINUX_64_BUILD_DIR}.tar.gz ${LINUX_64_BUILD_DIR}

tarball-freebsd64: ${DIST_DIR}/${FREEBSD_64_BUILD_DIR}.tar.gz
${DIST_DIR}/${FREEBSD_64_BUILD_DIR}.tar.gz: $(wildcard ${FREEBSD_64_BUILD_DIR_PATH} ${FREEBSD_64_BUILD_DIR_PATH}/*)
	cd ${DIST_DIR} ; tar cfz ${FREEBSD_64_BUILD_DIR}.tar.gz ${FREEBSD_64_BUILD_DIR}

tarball-macos64: ${DIST_DIR}/${MACOS_64_BUILD_DIR}.tar.gz
${DIST_DIR}/${MACOS_64_BUILD_DIR}.tar.gz: $(wildcard ${MACOS_64_BUILD_DIR_PATH} ${MACOS_64_BUILD_DIR_PATH}/*)
	cd ${DIST_DIR} ; tar cfz ${MACOS_64_BUILD_DIR}.tar.gz ${MACOS_64_BUILD_DIR}

zipball-win32: ${DIST_DIR}/${WIN_32_BUILD_DIR}.zip
${DIST_DIR}/${WIN_32_BUILD_DIR}.zip: $(wildcard ${WIN_32_BUILD_DIR_PATH} ${WIN_32_BUILD_DIR_PATH}/*)
	cd ${DIST_DIR} ; zip -qr ${WIN_32_BUILD_DIR}.zip ${WIN_32_BUILD_DIR}

zipball-win64: ${DIST_DIR}/${WIN_64_BUILD_DIR}.zip
${DIST_DIR}/${WIN_64_BUILD_DIR}.zip: $(wildcard ${WIN_64_BUILD_DIR_PATH} ${WIN_64_BUILD_DIR_PATH}/*)
	cd ${DIST_DIR} ; zip -qr ${WIN_64_BUILD_DIR}.zip ${WIN_64_BUILD_DIR}


### Clean up ###################################################################
#
# realclean = clean + distclean
#

# Remove whatever has been produced so far.
realclean: clean distclean

# Remove cross compiled output for various platforms and related tarballs.
distclean:
	rm -rf ${DIST_DIR}

# Remove output for the current platform.
clean:
	go clean ${SERVER_SRC_DIR_PATH}
	rm -rf ${BIN_BUILD_DIR_PATH}
	rm -rf ${UI_BUILD_DIR_PATH}
