.PHONY: build,compile,run

export CGO_ENABLED=0
export GO_APP=mmlicense
export GO_BIN=mmlicense

build: compile

compile:
	go build -o ./build/${GO_BIN} ${GO_APP}
	cp ./license.json ./build/

package:
	upx --brute ./build/${GO_BIN}

release:
	go build -ldflags="-s -w" -o ./build/${GO_BIN} ${GO_APP}
	upx --brute ./build/${GO_BIN}
	cp ./license.json ./build/

run: build
	cd ./build && ./${GO_BIN}