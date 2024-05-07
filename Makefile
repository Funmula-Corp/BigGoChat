export IGNORE_GO_WORK_IF_EXISTS=true

build: build-server build-package

build-package:
	cd server && $(MAKE) package

build-server:
	cd server && $(MAKE) build

run: run-server run-client

run-client:
	cd server && $(MAKE) run-client

run-server:
	cd server && $(MAKE) run-server

stop: stop-client stop-server

stop-client:
	cd server && $(MAKE) stop-client

stop-server:
	cd server && $(MAKE) stop-server

restart: restart-server restart-client

restart-server:
	cd server && $(MAKE) restart-server

restart-client:
	cd server && $(MAKE) restart-client
