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
	-pkill -INT -f hertz_gateway
	-pkill -f kitex_services
	-lsof -t -i :8870 | xargs kill
	-lsof -t -i :9870 | xargs kill
	-~/nacos/bin/shutdown.sh

