#
# @file
# Footle's makefile.
#
# Known to work with GNU make in GNU/Linux.
#
# Puts all build output into the "build/" directory.  Successful builds should
# produce *at least* an executable at "build/footle/bin/footle" and necessary
# HTML markup in "build/footle/ui/"
#
# Required tools:
# - Go >= 1.8: https://golang.org/doc/install
# - Go dep: go get -u github.com/golang/dep/cmd/dep
# - Node.js and npm: https://github.com/creationix/nvm
#
# Useful make commands:
# - make: Compiles Footle for the current platform.
# - make dist: Cross compiles for FreeBSD, Linux, Mac, and Windows.
# - make clean: Removes build for the current platform.
# - make distclean: Removes all cross compiled builds.
# - make realclean: Removes all types of builds.

BUILD_ROOT = ./build

# Build director for the current platform.
BUILD_DIR_PATH = ${BUILD_ROOT}/footle

# Build directories for all platforms.
LINUX_32_BUILD_DIR = footle-linux-32
LINUX_64_BUILD_DIR = footle-linux-64
FREEBSD_64_BUILD_DIR = footle-freebsd-64
MACOS_64_BUILD_DIR   = footle-mac-64
WIN_32_BUILD_DIR = footle-win-32
WIN_64_BUILD_DIR = footle-win-64
LINUX_32_BUILD_DIR_PATH = ${BUILD_ROOT}/${LINUX_32_BUILD_DIR}
LINUX_64_BUILD_DIR_PATH = ${BUILD_ROOT}/${LINUX_64_BUILD_DIR}
FREEBSD_64_BUILD_DIR_PATH = ${BUILD_ROOT}/${FREEBSD_64_BUILD_DIR}
MACOS_64_BUILD_DIR_PATH   = ${BUILD_ROOT}/${MACOS_64_BUILD_DIR}
WIN_32_BUILD_DIR_PATH = ${BUILD_ROOT}/${WIN_32_BUILD_DIR}
WIN_64_BUILD_DIR_PATH = ${BUILD_ROOT}/${WIN_64_BUILD_DIR}

# Location of Go code.
SERVER_SRC_DIR_PATH = ./src/server
GO_SRC_FILES = $(shell find ${SERVER_SRC_DIR_PATH} -name *.go -type f)
OUR_GO_PATH = $(shell pwd)

# Location of HTML/CSS/Javascript-based UI code.
UI_SRC_DIR_PATH = ./src/ui
UI_BUILD_DIR_PATH  = ${BUILD_DIR_PATH}/ui
UI_SRC_FILES = $(shell find ${UI_SRC_DIR_PATH} -type f)

UI_SRC_LIBS_DIR = ${UI_SRC_DIR_PATH}/libs
UI_BUILD_LIBS_DIR = ${UI_BUILD_DIR_PATH}/libs

UI_HTML_SRC_PATH = ${UI_SRC_DIR_PATH}/index.html
UI_HTML_BUILD_PATH = ${UI_BUILD_DIR_PATH}/index.html

UI_SASS_DIR = ${UI_SRC_DIR_PATH}/style/sass
UI_CSS_DIR = ${UI_BUILD_DIR_PATH}/style/css
UI_SASS_FILES = $(shell find ${UI_SASS_DIR} -type f)

UI_SCRIPT_SRC_DIR_PATH = ${UI_SRC_DIR_PATH}/scripts
UI_SCRIPT_BUILD_DIR_PATH = ${UI_BUILD_DIR_PATH}/scripts
UI_SCRIPT_SRC_FILES = $(shell find ${UI_SCRIPT_SRC_DIR_PATH} -type f)

UI_FONT_ORIG_DIR = ${UI_BUILD_DIR_PATH}/libs/bower_components/uikit/fonts
UI_FONT_REQUIRED_DIR = ${UI_BUILD_DIR_PATH}/style/fonts


# Compile Footle server's Go code and Footle UI's CSS/HTML/Javascript code.
all: Server ui


### Compile Footle server ######################################################
#
# Trivia: "server" as a target name does not work, but "Server" does.
#
Server: godeps ${BUILD_DIR_PATH}/bin/footle

# *Temporarily* set $GOPATH to Footle's top level dir so that Footle's Go
# packages (e.g. server/config) can be located during compilation.  This
# frees Footle from the Go workspace allowing it to reside anywhere
# in the file system.
${BUILD_DIR_PATH}/bin/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} go build -o ${BUILD_DIR_PATH}/bin/footle ${SERVER_SRC_DIR_PATH}

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


### Compile and prepare the UI #################################################
# Prepare HTML-based UI.
ui: ui-dev-dependencies ui-libs markup style script font

# Node.js packages.
ui-dev-dependencies: ${UI_SRC_DIR_PATH}/node_modules

${UI_SRC_DIR_PATH}/node_modules: ${UI_SRC_DIR_PATH}/package-lock.json
	cd ${UI_SRC_DIR_PATH}; \
	npm install

# Bower dependencies need to be installed in both UI src and UI build.  This
# is because some libs are needed during the build process.
ui-libs: ${UI_SRC_LIBS_DIR}/bower_components ${UI_BUILD_LIBS_DIR}/bower_components

${UI_SRC_LIBS_DIR}/bower_components: ${UI_SRC_LIBS_DIR}/bower.json
	cd ${UI_SRC_LIBS_DIR}; \
	npx bower install

${UI_BUILD_LIBS_DIR}/bower_components: ${UI_BUILD_LIBS_DIR}/bower.json
	cd ${UI_BUILD_LIBS_DIR}; \
	npx bower install

${UI_BUILD_LIBS_DIR}/bower.json: ${UI_SRC_LIBS_DIR}/bower.json
	mkdir -p ${UI_BUILD_LIBS_DIR}
	cp $< $@

# Prepare CSS whenever *any* Sass file changes.
style: ${UI_CSS_DIR}/ui.css
${UI_CSS_DIR}/ui.css: ${UI_SASS_FILES}
	mkdir -p ${UI_CSS_DIR}
	npx node-sass --indented-syntax --source-map true ${UI_SASS_DIR}/ui.sass $@

# Copy index.html
markup: ${UI_HTML_BUILD_PATH}
${UI_HTML_BUILD_PATH}: ${UI_HTML_SRC_PATH}
	cp $< $@

# Copy Javascript whenever *any* script file changes.
script: ${UI_SCRIPT_BUILD_DIR_PATH} 
${UI_SCRIPT_BUILD_DIR_PATH}: ${UI_SCRIPT_SRC_DIR_PATH}
	cp -r ${UI_SCRIPT_SRC_DIR_PATH}/. $@
	touch $@

# Relocate uikit's "fonts" directory.  Otherwise it won't be found by the
# web browser.
font: ${UI_FONT_REQUIRED_DIR}
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
test:
	GOPATH=${OUR_GO_PATH} go test ${SERVER_SRC_DIR_PATH}/core/... ${SERVER_SRC_DIR_PATH}/dbgp/... ${SERVER_SRC_DIR_PATH}/http/...


### Cross compile and prepare tarballs ########################################
#
# Prepare distributions for various platforms.
#
# Preparation formula: Cross compile Go code, copy UI code, create tar/zip
# archive.
#
dist: cross-compile ui-copy tarball

cross-compile: linux32 linux64 freebsd64 macos64 win32 win64

linux32: ${LINUX_32_BUILD_DIR_PATH}/bin/footle
${LINUX_32_BUILD_DIR_PATH}/bin/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=linux GOARCH=386 go build -o ${LINUX_32_BUILD_DIR_PATH}/bin/footle ${SERVER_SRC_DIR_PATH}

linux64: ${LINUX_64_BUILD_DIR_PATH}/bin/footle
${LINUX_64_BUILD_DIR_PATH}/bin/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=linux GOARCH=amd64 go build -o ${LINUX_64_BUILD_DIR_PATH}/bin/footle ${SERVER_SRC_DIR_PATH}

freebsd64: ${FREEBSD_64_BUILD_DIR_PATH}/bin/footle
${FREEBSD_64_BUILD_DIR_PATH}/bin/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=freebsd GOARCH=amd64 go build -o ${FREEBSD_64_BUILD_DIR_PATH}/bin/footle ${SERVER_SRC_DIR_PATH}

macos64: ${MACOS_64_BUILD_DIR_PATH}/bin/footle
${MACOS_64_BUILD_DIR_PATH}/bin/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=darwin GOARCH=amd64 go build -o ${MACOS_64_BUILD_DIR_PATH}/bin/footle ${SERVER_SRC_DIR_PATH}

win32: ${WIN_32_BUILD_DIR_PATH}/bin/footle
${WIN_32_BUILD_DIR_PATH}/bin/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=windows GOARCH=386 go build -o ${WIN_32_BUILD_DIR_PATH}/bin/footle ${SERVER_SRC_DIR_PATH}

win64: ${WIN_64_BUILD_DIR_PATH}/bin/footle
${WIN_64_BUILD_DIR_PATH}/bin/footle: ${GO_SRC_FILES}
	GOPATH=${OUR_GO_PATH} GOOS=windows GOARCH=amd64 go build -o ${WIN_64_BUILD_DIR_PATH}/bin/footle ${SERVER_SRC_DIR_PATH}

# Copy the same UI code inside every distribution everytime *any* UI source
# file changes.  Do not forget to update the timestamp of the *copied* directory
# to avoid further unnecessary copying.
ui-copy: ui ui-copy-linux32 ui-copy-linux64 ui-copy-freebsd64 ui-copy-macos64 ui-copy-win32 ui-copy-win64

ui-copy-linux32: ${LINUX_32_BUILD_DIR_PATH}/ui
${LINUX_32_BUILD_DIR_PATH}/ui: ${UI_SRC_FILES}
	mkdir -p ${LINUX_32_BUILD_DIR_PATH}
	cp -r ${UI_BUILD_DIR_PATH} ${LINUX_32_BUILD_DIR_PATH}
	touch ${LINUX_32_BUILD_DIR_PATH}/ui

ui-copy-linux64: ${LINUX_64_BUILD_DIR_PATH}/ui
${LINUX_64_BUILD_DIR_PATH}/ui: ${UI_SRC_FILES}
	mkdir -p ${LINUX_64_BUILD_DIR_PATH}
	cp -r ${UI_BUILD_DIR_PATH} ${LINUX_64_BUILD_DIR_PATH}
	touch ${LINUX_64_BUILD_DIR_PATH}/ui

ui-copy-freebsd64: ${FREEBSD_64_BUILD_DIR_PATH}/ui
${FREEBSD_64_BUILD_DIR_PATH}/ui: ${UI_SRC_FILES}
	mkdir -p ${FREEBSD_64_BUILD_DIR_PATH}
	cp -r ${UI_BUILD_DIR_PATH} ${FREEBSD_64_BUILD_DIR_PATH}
	touch ${FREEBSD_64_BUILD_DIR_PATH}/ui

ui-copy-macos64: ${MACOS_64_BUILD_DIR_PATH}/ui
${MACOS_64_BUILD_DIR_PATH}/ui: ${UI_SRC_FILES}
	mkdir -p ${MACOS_64_BUILD_DIR_PATH}
	cp -r ${UI_BUILD_DIR_PATH} ${MACOS_64_BUILD_DIR_PATH}
	touch ${MACOS_64_BUILD_DIR_PATH}/ui

ui-copy-win32: ${WIN_32_BUILD_DIR_PATH}/ui
${WIN_32_BUILD_DIR_PATH}/ui: ${UI_SRC_FILES}
	mkdir -p ${WIN_32_BUILD_DIR_PATH}
	cp -r ${UI_BUILD_DIR_PATH} ${WIN_32_BUILD_DIR_PATH}
	touch ${WIN_32_BUILD_DIR_PATH}/ui

ui-copy-win64: ${WIN_64_BUILD_DIR_PATH}/ui
${WIN_64_BUILD_DIR_PATH}/ui: ${UI_SRC_FILES}
	mkdir -p ${WIN_64_BUILD_DIR_PATH}
	cp -r ${UI_BUILD_DIR_PATH} ${WIN_64_BUILD_DIR_PATH}
	touch ${WIN_64_BUILD_DIR_PATH}/ui

tarball: tarball-linux32 tarball-linux64 tarball-freebsd64 tarball-macos64 zipball-win32 zipball-win64

tarball-linux32: ${BUILD_ROOT}/${LINUX_32_BUILD_DIR}.tar.gz
${BUILD_ROOT}/${LINUX_32_BUILD_DIR}.tar.gz: ${LINUX_32_BUILD_DIR_PATH}/bin ${LINUX_32_BUILD_DIR_PATH}/ui
	cd ${BUILD_ROOT} ; tar cfz ${LINUX_32_BUILD_DIR}.tar.gz ${LINUX_32_BUILD_DIR}

tarball-linux64: ${BUILD_ROOT}/${LINUX_64_BUILD_DIR}.tar.gz
${BUILD_ROOT}/${LINUX_64_BUILD_DIR}.tar.gz: ${LINUX_64_BUILD_DIR_PATH}/bin ${LINUX_64_BUILD_DIR_PATH}/ui
	cd ${BUILD_ROOT} ; tar cfz ${LINUX_64_BUILD_DIR}.tar.gz ${LINUX_64_BUILD_DIR}

tarball-freebsd64: ${BUILD_ROOT}/${FREEBSD_64_BUILD_DIR}.tar.gz
${BUILD_ROOT}/${FREEBSD_64_BUILD_DIR}.tar.gz: ${FREEBSD_64_BUILD_DIR_PATH}/bin ${FREEBSD_64_BUILD_DIR_PATH}/ui
	cd ${BUILD_ROOT} ; tar cfz ${FREEBSD_64_BUILD_DIR}.tar.gz ${FREEBSD_64_BUILD_DIR}

tarball-macos64: ${BUILD_ROOT}/${MACOS_64_BUILD_DIR}.tar.gz
${BUILD_ROOT}/${MACOS_64_BUILD_DIR}.tar.gz: ${MACOS_64_BUILD_DIR_PATH}/bin ${MACOS_64_BUILD_DIR_PATH}/ui
	cd ${BUILD_ROOT} ; tar cfz ${MACOS_64_BUILD_DIR}.tar.gz ${MACOS_64_BUILD_DIR}

zipball-win32: ${BUILD_ROOT}/${WIN_32_BUILD_DIR}.zip
${BUILD_ROOT}/${WIN_32_BUILD_DIR}.zip: ${WIN_32_BUILD_DIR_PATH}/bin ${WIN_32_BUILD_DIR_PATH}/ui
	cd ${BUILD_ROOT} ; zip -qr ${WIN_32_BUILD_DIR}.zip ${WIN_32_BUILD_DIR}

zipball-win64: ${BUILD_ROOT}/${WIN_64_BUILD_DIR}.zip
${BUILD_ROOT}/${WIN_64_BUILD_DIR}.zip: ${WIN_64_BUILD_DIR_PATH}/bin ${WIN_64_BUILD_DIR_PATH}/ui
	cd ${BUILD_ROOT} ; zip -qr ${WIN_64_BUILD_DIR}.zip ${WIN_64_BUILD_DIR}


### Clean up ###################################################################
#
# realclean = clean + distclean
#

# Remove whatever has been produced so far.
realclean: clean distclean

# Remove cross compiled output for various platforms and related tarballs.
distclean:
	cd ${BUILD_ROOT} ; \
	rm -rf ${LINUX_32_BUILD_DIR} ${LINUX_32_BUILD_DIR}.tar.gz \
					${LINUX_64_BUILD_DIR} ${LINUX_64_BUILD_DIR}.tar.gz \
					${FREEBSD_64_BUILD_DIR} ${FREEBSD_64_BUILD_DIR}.tar.gz \
					${MACOS_64_BUILD_DIR} ${MACOS_64_BUILD_DIR}.tar.gz \
					${WIN_32_BUILD_DIR} ${WIN_32_BUILD_DIR}.zip \
					${WIN_64_BUILD_DIR} ${WIN_64_BUILD_DIR}.zip

# Remove output for the current platform.
clean:
	go clean ${SERVER_SRC_DIR_PATH}
	rm -rf ${BUILD_DIR_PATH}
