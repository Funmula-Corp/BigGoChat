export IGNORE_GO_WORK_IF_EXISTS=true

build: build-server build-package

build-package:
	cd server && $(MAKE) package

build-server:
	cd server && $(MAKE) build
