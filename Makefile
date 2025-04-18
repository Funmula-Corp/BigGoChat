export IGNORE_GO_WORK_IF_EXISTS=true
export GOPRIVATE=git.biggo.com

export MM_SERVICEENVIRONMENT=production
#export MM_LICENSE=eyJpZCI6ImM3YmhuQXhwUVBMSEZlQkpaOVE2ZEJNN1ViIiwiaXNzdWVkX2F0IjoxNzE1MjQ4MTA5NzMxLCJzdGFydHNfYXQiOjE3MTUxODQwMDAwMDAsImV4cGlyZXNfYXQiOjE3NDY3MjAwMDAwMDAsImN1c3RvbWVyIjp7ImlkIjoiYmlnZ29fZGV2X3RlYW0iLCJuYW1lIjoiQiFnR28gRGV2ZWxvcG1lbnQiLCJlbWFpbCI6InJvb3RAbG9jYWxob3N0IiwiY29tcGFueSI6IkZ1bm11bGEifSwiZmVhdHVyZXMiOnsidXNlcnMiOjIsImxkYXAiOmZhbHNlLCJsZGFwX2dyb3VwcyI6ZmFsc2UsIm1mYSI6ZmFsc2UsImdvb2dsZV9vYXV0aCI6ZmFsc2UsIm9mZmljZTM2NV9vYXV0aCI6ZmFsc2UsIm9wZW5pZCI6ZmFsc2UsImNvbXBsaWFuY2UiOmZhbHNlLCJjbHVzdGVyIjp0cnVlLCJtZXRyaWNzIjpmYWxzZSwibWhwbnMiOmZhbHNlLCJzYW1sIjpmYWxzZSwiZWxhc3RpY19zZWFyY2giOnRydWUsImFubm91bmNlbWVudCI6ZmFsc2UsInRoZW1lX21hbmFnZW1lbnQiOmZhbHNlLCJlbWFpbF9ub3RpZmljYXRpb25fY29udGVudHMiOmZhbHNlLCJkYXRhX3JldGVudGlvbiI6ZmFsc2UsIm1lc3NhZ2VfZXhwb3J0IjpmYWxzZSwiY3VzdG9tX3Blcm1pc3Npb25zX3NjaGVtZXMiOmZhbHNlLCJjdXN0b21fdGVybXNfb2Zfc2VydmljZSI6ZmFsc2UsImd1ZXN0X2FjY291bnRzIjpmYWxzZSwiZ3Vlc3RfYWNjb3VudHNfcGVybWlzc2lvbnMiOmZhbHNlLCJpZF9sb2FkZWQiOmZhbHNlLCJsb2NrX3RlYW1tYXRlX25hbWVfZGlzcGxheSI6ZmFsc2UsImVudGVycHJpc2VfcGx1Z2lucyI6ZmFsc2UsImFkdmFuY2VkX2xvZ2dpbmciOmZhbHNlLCJjbG91ZCI6ZmFsc2UsInNoYXJlZF9jaGFubmVscyI6ZmFsc2UsInJlbW90ZV9jbHVzdGVyX3NlcnZpY2UiOnRydWUsIm91dGdvaW5nX29hdXRoX2Nvbm5lY3Rpb25zIjpmYWxzZSwiZnV0dXJlX2ZlYXR1cmVzIjpmYWxzZX0sInNrdV9uYW1lIjoiQiFnR28gQ2hhdCBMaWNlbnNlIiwic2t1X3Nob3J0X25hbWUiOiJCIWdHbyBDaGF0IExpY2Vuc2UiLCJpc190cmlhbCI6ZmFsc2UsImlzX2dvdl9za3UiOnRydWUsInNpZ251cF9qd3QiOm51bGx9Xk6wg66C2MeOvgTYCTkln0ldmClgW7ZNqobj1+aGRdKToAm6rOS4OKkcaU1gskA4klTN1jFdNDFtMIRopS9ySfkvr7wByNAVAykm/Gbr+nAs3x0VIYi1T3P/I2Ez21/O5RjVp7f9+tEowlY9D91iSznkuX9FsaqBaWnzuN9GWvYkXLQtu0O8l0HwV2P2V3fzRLBPgocXVZ/F7SSCUYXNhrtwBBR+ZDfvw1nkjKisEtBvR2OXKj7Nlgoxm30Zn11aUdcQLd0Iht2HReU/9G5OAQYXoWyjW/DFALO5PCCYuekPXzMLwqVpoFg0yFkRmxMc9hkv81pUoDHYg1kaFaxHIQ==

DOCKER_HOST:=docker.dev.cloud.biggo.com
DOCKER_PROJECT:=test
DOCKER_IMAGE:=mattermost
DOCKER_TAG:=cluster-dev

build: build-server build-package

build-package:
	cd server && $(MAKE) package

build-docker:
	cd server/build && \
	cp -r ../dist packages && \
	docker build -t ${DOCKER_HOST}/${DOCKER_PROJECT}/${DOCKER_IMAGE}:${DOCKER_TAG} .

build-server:
	cd server && $(MAKE) build

build-webapp:
	cd webapp && $(MAKE) dist

clean-docker:
	cd server && $(MAKE) clean-docker

debug-headless:
	cd server && $(MAKE) debug-server-headless

gen-all:
	cd server && $(MAKE) pluginapi plugin-mocks generated gen-serialized

run: run-server run-client

run-client:
	cd server && $(MAKE) run-client

run-server:
	cd server && $(MAKE) run-server

stop: stop-client stop-server

stop-client:
	cd server && $(MAKE) stop-client

stop-docker:
	cd server && $(MAKE) stop-docker

stop-server:
	cd server && $(MAKE) stop-server

restart: restart-server restart-client

restart-server:
	cd server && $(MAKE) restart-server

restart-client:
	cd server && $(MAKE) restart-client

run-docker:
	docker run -it --rm -v ./server/config/config.json:/mattermost/config/config.json ${DOCKER_HOST}/${DOCKER_PROJECT}/${DOCKER_IMAGE}:${DOCKER_TAG}

build-full: build-server build-package build-docker publish-docker

publish-docker:
	docker push ${DOCKER_HOST}/${DOCKER_PROJECT}/${DOCKER_IMAGE}:${DOCKER_TAG}

build-drone-docker-image:
	docker build --tag docker.cloud.biggo.com/droneio/mattermost-funmula-build:latest --file ./Dockerfile.drone .
	docker push docker.cloud.biggo.com/droneio/mattermost-funmula-build:latest

# https://git.biggo.com/Funmula/BigGoChat/wiki/Write-e2e-test-for-frontend-update
run-cypress:
	@cd ./e2e-tests/cypress && npx cypress open
