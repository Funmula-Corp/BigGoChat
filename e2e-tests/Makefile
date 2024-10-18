export MM_LICENSE=eyJpZCI6InNrcTA4bE5ZUDhmQnZzSmZDcVJ4amJYckVxIiwiaXNzdWVkX2F0IjoxNzI0OTAzMDY3OTEyLCJzdGFydHNfYXQiOjE3MjQ4NjA4MDAwMDAsImV4cGlyZXNfYXQiOjE3NTYzOTY4MDAwMDAsImN1c3RvbWVyIjp7ImlkIjoiYmlnZ29fZGV2X3RlYW0iLCJuYW1lIjoiQiFnR28gRGV2ZWxvcG1lbnQiLCJlbWFpbCI6InJvb3RAbG9jYWxob3N0IiwiY29tcGFueSI6IkZ1bm11bGEifSwiZmVhdHVyZXMiOnsidXNlcnMiOjEwMDAwMDAsImxkYXAiOnRydWUsImxkYXBfZ3JvdXBzIjp0cnVlLCJtZmEiOnRydWUsImdvb2dsZV9vYXV0aCI6dHJ1ZSwib2ZmaWNlMzY1X29hdXRoIjp0cnVlLCJvcGVuaWQiOmZhbHNlLCJjb21wbGlhbmNlIjp0cnVlLCJjbHVzdGVyIjp0cnVlLCJtZXRyaWNzIjp0cnVlLCJtaHBucyI6dHJ1ZSwic2FtbCI6dHJ1ZSwiZWxhc3RpY19zZWFyY2giOnRydWUsImFubm91bmNlbWVudCI6dHJ1ZSwidGhlbWVfbWFuYWdlbWVudCI6dHJ1ZSwiZW1haWxfbm90aWZpY2F0aW9uX2NvbnRlbnRzIjp0cnVlLCJkYXRhX3JldGVudGlvbiI6dHJ1ZSwibWVzc2FnZV9leHBvcnQiOnRydWUsImN1c3RvbV9wZXJtaXNzaW9uc19zY2hlbWVzIjp0cnVlLCJjdXN0b21fdGVybXNfb2Zfc2VydmljZSI6dHJ1ZSwiZ3Vlc3RfYWNjb3VudHMiOnRydWUsImd1ZXN0X2FjY291bnRzX3Blcm1pc3Npb25zIjp0cnVlLCJpZF9sb2FkZWQiOnRydWUsImxvY2tfdGVhbW1hdGVfbmFtZV9kaXNwbGF5Ijp0cnVlLCJlbnRlcnByaXNlX3BsdWdpbnMiOnRydWUsImFkdmFuY2VkX2xvZ2dpbmciOnRydWUsImNsb3VkIjpmYWxzZSwic2hhcmVkX2NoYW5uZWxzIjp0cnVlLCJyZW1vdGVfY2x1c3Rlcl9zZXJ2aWNlIjp0cnVlLCJvdXRnb2luZ19vYXV0aF9jb25uZWN0aW9ucyI6dHJ1ZSwiZnV0dXJlX2ZlYXR1cmVzIjp0cnVlfSwic2t1X25hbWUiOiJCIWdHbyBDaGF0IExpY2Vuc2UiLCJza3Vfc2hvcnRfbmFtZSI6IkIhZ0dvIENoYXQgTGljZW5zZSIsImlzX3RyaWFsIjpmYWxzZSwiaXNfZ292X3NrdSI6ZmFsc2UsInNpZ251cF9qd3QiOm51bGx9aiHv2mEnLZFFX7dnuS1ogtk0il/XIbi7JtByAR5vyIM6PnHIoVBRqN1yH//MG5WXSRRCkUfoZ0sCBJiaUiA/nsMnmCvsE439rhcm9HaQ1sdNqfXCFj+P4WIn+E0IiGmNgSF6wKnD6R35bBG86prjYQ7+jmGfsSl5Ez01VolRepV6E5CD5F6w+58ewuNqeS4wa0qXtaWVZX6U/WdXWvb2WbZTq+oju2YwtDG/3XU43L3XzZovWwhQMY3PR/b0p7ICaDho559RvqOjYX+UuTZAAwLNsiowj75OXASXtzmEU/0V2zXTPpAzIDPsZAHfgFzApJ4H1lBVo0AkPjdXCsau6g==
export TEST_FILTER=''

SHELL := /bin/bash

.PHONY: all run stop clean
all: run
run: generate-server start-server run-test
stop: stop-server stop-dashboard clean
clean:
	rm -fv .ci/server.yml
	rm -fv .ci/.env.{server,dashboard,cypress,playwright}

.PHONY: generate-server start-server run-test stop-server restart-server
generate-server:
	bash ./.ci/server.generate.sh
start-server: generate-server
	bash ./.ci/server.start.sh
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
