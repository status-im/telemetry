# Telemetry

Opt-in message reliability metrics service.

# Development

You need to setup a postgres db as such:
1) Create a telemetry user with password newPassword
2) Create a db telemetry
3) Create a db telemetry_test

Then you can run the server with:
```
go run cmd/server/main.go -data-source-name postgres://telemetry:newPassword@127.0.0.1:5432/telemetry
```

Finally, to run the test:
```
make test
```

# Continuous Integration

Builds of Docker images are created with our [Jenkins CI](https://ci.infra.status.im/job/telemetry/job/docker/) which push a [`statusteam/telemetry`](https://hub.docker.com/r/statusteam/telemetry) Docker image.

The host is managed in [`infra-misc`](https://github.com/status-im/infra-misc/blob/master/ansible/roles/telemetry) repository.
