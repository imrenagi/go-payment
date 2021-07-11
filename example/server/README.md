Example
===

To start using this module, you can try the example [server.go](/server.go)

This example is using sqlite as the datastore. If you want to change the database,
please read gorm docs for connecting to database [here](https://gorm.io/docs/connecting_to_the_database.html)

To run the application, simply use:

```console
$ go run example/server/server.go
```

> :heavy_exclamation_mark: If you want to accept payment callback from the payment gateway on your local computer for development purpose, consider to use [ngrok.io](https://ngrok.io) to expose your localhost to the internet and update the callback base URL in all url environment variables.