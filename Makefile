VERSION     = 1.0.1
BUILD_DATE  = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: build

build:
	sudo docker build \
		--rm --compress \
		--build-arg VERSION="$(VERSION)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--tag imagelist/smartctl_exporter:latest \
		--tag imagelist/smartctl_exporter:$(VERSION) \
		.

build_go:
	go build -o smartctl_exporter -ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}"

push:
	docker push imagelist/smartctl_exporter:latest
	docker push imagelist/smartctl_exporter:$(VERSION)
