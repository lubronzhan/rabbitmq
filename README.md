# rabbitmq
Learning


1. Basic receive and sender. Open 3 terminals

```sh
go run 1-receive.go
```

```sh
go run 1-receive.go
```

```sh
go run 1-send.go
```

2. Long running task, wait for execute to finish before ack
```sh
go run 2-worker.go
```

```sh
go run 2-worker.go
```

```sh
go run 2-sendworker.go message.
go run 2-sendworker.go message....
go run 2-sendworker.go message........
```

3. Fanout to all queue

```sh
go run 3-receive_log.go
```

```sh
go run 3-receive_log.go
```

```sh
go run 3-emit_log.go message1
go run 3-emit_log.go message2
go run 3-emit_log.go message3
```

4. Direct to specified server

```sh
go run 3-receive_log_direct.go info
```

```sh
go run 3-receive_log_direct.go fatal
```

```sh
go run 3-emit_log_direct.go info message1
go run 3-emit_log_direct.go fatal message2
go run 3-emit_log_direct.go info message3
```

5. Direct to server match the regex routing_key

```sh
go run 4-receive_log_topic.go "*.critical"
```

```sh
go run 4-receive_log_topic.go "kern.*"
```

```sh
go run 4-emit_log_topic.go "asdf-asdf.critical" body1
// won't match
go run 4-emit_log_topic.go ".asdf-asdf.critical" body
go run 4-emit_log_topic.go "kern.asdf" message3
```

6. grpc

```sh
go run 6-rpc_server.go
```

```sh
go run 6-rpc_server.go
```

```sh
go run 6-rpc_client.go 30
go run 6-rpc_client.go 20
```
