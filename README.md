# How to run

## The simplest way

Make sure you have installed docker-compose. And run:

```
docker-compose up
```

Wait a sec for the initializations and you'll see a log like:

```
boxpractice-boxpractice-1  | {"time":"2022-09-05T09:44:27.593Z","logger":"setup","caller":"boxpractice/main.go:63","msg":"start server","v":0}
```

Now you can access `http://localhost:8081` to see the API docs and also can test the API using the Swagger UI.
