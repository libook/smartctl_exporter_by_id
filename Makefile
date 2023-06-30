VERSION     = 1.0.2
BUILD_DATE  = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: build

build:
	sudo docker build \
		--rm --compress \
		--build-arg VERSION="$(VERSION)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--tag libook/smartctl_exporter_by_id:latest \
		--tag libook/smartctl_exporter_by_id:$(VERSION) \
		.

build_go:
	go build -o smartctl_exporter_by_id -ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}"

push:
	docker push libook/smartctl_exporter_by_id:latest
	docker push libook/smartctl_exporter_by_id:$(VERSION)
