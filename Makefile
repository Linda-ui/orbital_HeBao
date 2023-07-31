.SILENT:

.PHONY: nacos
nacos:
	~/nacos/bin/startup.sh -m standalone

.PHONY: gateway
gateway: nacos
	go run ./hertz_gateway

.PHONY: services
services:
	./scripts/server_startup.sh

.PHONY: stop
stop:
	-pkill -f INT hertz_gateway
	-pkill -f kitex_services
	-~/nacos/bin/shutdown.sh

