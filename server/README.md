# How to run the server side services

First, make sure the RabbitMQ broker is running. You can do that by executing:
```docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management```

Then, start the orchestrator (first example), as well as its backup (second example). The orcherstrator address is up for you to decide, but you must provide the address of its backup as well. Run the following commands (on two different terminals) in the `server/orchestrator` directory:

```go run orchestrator.go localhost:8080 localhost:8081``` and ```go run orchestrator.go localhost:8081 localhost:8080```

Lastly, start as many servers as your heart desires, but don't forget to tell it both of the orchestrator addresses! You can do it by executing, in different terminals, in the `server` directory the command. You also need to specify the database name to be used by the server in the last argument. The command would be:
```go run main.go localhost:8080 localhost:8081 <DB_NAME>```

Then, to simulate a client-side message sender, you can run, in the `server/mock/client` directory, the following command:
```go run mockClient.go```

Important to note, the minimum number of quorum participants is determined by the number of active servers (`NumberOfServers/2 + 1`), so there needs to be at least **two** servers running in order for a quorum to happen.
