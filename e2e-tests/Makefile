#export MM_LICENSE=eyJpZCI6ImM3YmhuQXhwUVBMSEZlQkpaOVE2ZEJNN1ViIiwiaXNzdWVkX2F0IjoxNzE1MjQ4MTA5NzMxLCJzdGFydHNfYXQiOjE3MTUxODQwMDAwMDAsImV4cGlyZXNfYXQiOjE3NDY3MjAwMDAwMDAsImN1c3RvbWVyIjp7ImlkIjoiYmlnZ29fZGV2X3RlYW0iLCJuYW1lIjoiQiFnR28gRGV2ZWxvcG1lbnQiLCJlbWFpbCI6InJvb3RAbG9jYWxob3N0IiwiY29tcGFueSI6IkZ1bm11bGEifSwiZmVhdHVyZXMiOnsidXNlcnMiOjIsImxkYXAiOmZhbHNlLCJsZGFwX2dyb3VwcyI6ZmFsc2UsIm1mYSI6ZmFsc2UsImdvb2dsZV9vYXV0aCI6ZmFsc2UsIm9mZmljZTM2NV9vYXV0aCI6ZmFsc2UsIm9wZW5pZCI6ZmFsc2UsImNvbXBsaWFuY2UiOmZhbHNlLCJjbHVzdGVyIjp0cnVlLCJtZXRyaWNzIjpmYWxzZSwibWhwbnMiOmZhbHNlLCJzYW1sIjpmYWxzZSwiZWxhc3RpY19zZWFyY2giOnRydWUsImFubm91bmNlbWVudCI6ZmFsc2UsInRoZW1lX21hbmFnZW1lbnQiOmZhbHNlLCJlbWFpbF9ub3RpZmljYXRpb25fY29udGVudHMiOmZhbHNlLCJkYXRhX3JldGVudGlvbiI6ZmFsc2UsIm1lc3NhZ2VfZXhwb3J0IjpmYWxzZSwiY3VzdG9tX3Blcm1pc3Npb25zX3NjaGVtZXMiOmZhbHNlLCJjdXN0b21fdGVybXNfb2Zfc2VydmljZSI6ZmFsc2UsImd1ZXN0X2FjY291bnRzIjpmYWxzZSwiZ3Vlc3RfYWNjb3VudHNfcGVybWlzc2lvbnMiOmZhbHNlLCJpZF9sb2FkZWQiOmZhbHNlLCJsb2NrX3RlYW1tYXRlX25hbWVfZGlzcGxheSI6ZmFsc2UsImVudGVycHJpc2VfcGx1Z2lucyI6ZmFsc2UsImFkdmFuY2VkX2xvZ2dpbmciOmZhbHNlLCJjbG91ZCI6ZmFsc2UsInNoYXJlZF9jaGFubmVscyI6ZmFsc2UsInJlbW90ZV9jbHVzdGVyX3NlcnZpY2UiOnRydWUsIm91dGdvaW5nX29hdXRoX2Nvbm5lY3Rpb25zIjpmYWxzZSwiZnV0dXJlX2ZlYXR1cmVzIjpmYWxzZX0sInNrdV9uYW1lIjoiQiFnR28gQ2hhdCBMaWNlbnNlIiwic2t1X3Nob3J0X25hbWUiOiJCIWdHbyBDaGF0IExpY2Vuc2UiLCJpc190cmlhbCI6ZmFsc2UsImlzX2dvdl9za3UiOnRydWUsInNpZ251cF9qd3QiOm51bGx9Xk6wg66C2MeOvgTYCTkln0ldmClgW7ZNqobj1+aGRdKToAm6rOS4OKkcaU1gskA4klTN1jFdNDFtMIRopS9ySfkvr7wByNAVAykm/Gbr+nAs3x0VIYi1T3P/I2Ez21/O5RjVp7f9+tEowlY9D91iSznkuX9FsaqBaWnzuN9GWvYkXLQtu0O8l0HwV2P2V3fzRLBPgocXVZ/F7SSCUYXNhrtwBBR+ZDfvw1nkjKisEtBvR2OXKj7Nlgoxm30Zn11aUdcQLd0Iht2HReU/9G5OAQYXoWyjW/DFALO5PCCYuekPXzMLwqVpoFg0yFkRmxMc9hkv81pUoDHYg1kaFaxHIQ==

SHELL := /bin/bash

.PHONY: all run stop clean
all: run
run: cloud-init generate-server start-server prepare-server run-test cloud-teardown
stop: stop-server stop-dashboard clean
clean:
	rm -fv .ci/server.yml
	rm -fv .ci/.env*

.PHONY: generate-server start-server prepare-server run-test stop-server restart-server
generate-server:
	bash ./.ci/server.generate.sh
start-server: generate-server
	bash ./.ci/server.start.sh
prepare-server:
	bash ./.ci/server.prepare.sh
run-test:
	bash ./.ci/server.run_test.sh
stop-server: generate-server
	bash ./.ci/server.stop.sh
restart-server: stop-server start-server

.PHONY: start-dashboard generate-test-cycle stop-dashboard
start-dashboard:
	bash ./.ci/dashboard.start.sh
generate-test-cycle: requirecmd-node generate-server
	bash ./.ci/dashboard.generate_test_cycle.sh
stop-dashboard:
	bash ./.ci/dashboard.stop.sh

.PHONY: cloud-init cloud-teardown
cloud-init: requirecmd-jq requirecmd-curl
	bash ./.ci/server.cloud_init.sh
cloud-teardown:
	bash ./.ci/server.cloud_teardown.sh

.PHONY: publish-report
publish-report: requirecmd-node
	bash ./.ci/report.publish.sh

.PHONY: fmt-node fmt-shell fmt
requirecmd-%:
	@which "$(*)" >/dev/null || { echo "Error, missing required CLI tool: $(*). Aborting." >&2; exit 1; }
fmt-node: requirecmd-npx
	# Formats yaml files
	npx prettier ./.ci "!./.ci/dashboard" --write --cache
fmt-shell: requirecmd-shfmt requirecmd-shellcheck
	shfmt -w -s -i 2 ./.ci/*.sh         # Install with https://webinstall.dev/shellcheck/
	shellcheck ./.ci/*.sh ./.ci/.e2erc* # Install with https://webinstall.dev/shfmt/
fmt: fmt-node fmt-shell
