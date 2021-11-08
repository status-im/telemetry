# telemetry

## Dev setup

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
