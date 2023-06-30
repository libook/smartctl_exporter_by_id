# smartctl_exporter_by_id

[![Docker Image Size](https://badgen.net/docker/size/libook/smartctl_exporter_by_id?icon=docker&label=image%20size)](https://hub.docker.com/r/libook/smartctl_exporter_by_id)

Prometheus exporter for [smartmontools](https://www.smartmontools.org/) to export the S.M.A.R.T. attributes.

This exporter uses device IDs which created in `/dev/disk/by-id` as `device` in case labels(like `/dev/sda`) changed after hardware modification.

This project was modified from https://github.com/STI26/smartctl_exporter which is using label as device name. Thank the earlier work by [ctrysbita](https://github.com/ctrysbita) and [STI26](https://github.com/STI26).

## Deployment

```sh
docker run --detach --privileged -p 9111:9111 --name smartctl_exporter_by_id -v /dev:/dev:ro libook/smartctl_exporter_by_id:latest
# or via podman
podman run --detach --privileged -p 9111:9111 --name smartctl-exporter-by-id -v /dev:/dev:ro docker.io/libook/smartctl_exporter_by_id:latest
```

### Parameters
- `--addr` Listen address. Default to listen on all network interface.
- `--port` Listen port. Default to `9111`.
- `--path` Matrics path. Default to `/metrics`.
- `--disable-auth` Set to `true` to disable Basic Authentication.
- `--user` Username of Basic Authentication. Default to `admin`.
- `--pass` Password of Basic Authentication. Default to `admin`.
- `--version` Show version.

Metrics will be available at http://localhost:9111/metrics

## Grafana Dashboard

![](https://github.com/libook/smartctl_exporter_by_id/assets/3395610/e389adfb-6e14-430b-b8db-426c328aefb4)
