# How to run the server side services

First, make sure the RabbitMQ broker is running. You can do that by executing:
```docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management```

Then, start the orchestrator as well as it's backup. The orcherstrator address and port are hardcoded, so make sure you use port `8080` for the original orchestrator and port `8081` for it's backup. Run the following commands (on two different terminals) in the `server/orchestrator` directory:

```go run orchestrator.go 8080``` and ```go run orchestrator.go 8081```

Lastly, start as many servers as your heart desires! You can do it by executing, in different terminals, in the `server` directory the command:
```go run main.go```

Then, to simulate a client-side message sender, you can run, in the `server/mock/client` directory, the following command:
```go run mockClient.go```

Important to note, the minimum number of quorum participants is determined by the number of active servers (`NumberOfServers/2 + 1`), so there needs to be at least **two** servers running in order for a quorum to happen.