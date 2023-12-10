# How to run the client-side services

In case you wish to be connected to the server, please start the RabbitMQ broker by running:
```docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management```

Then, to start the client backend services please execute the following command in the directory `client/backend`:
```go run main.go```

Lastly, to start the frontend services execute in the directory `client/frontend`:
```npm i``` then ```npm run dev```
