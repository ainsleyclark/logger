# Logger

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Test](https://github.com/krang-backlink/logger/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/krang-backlink/logger/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/krang-backlink/logger/branch/master/graph/badge.svg?token=HXEZIJ2NUY)](https://codecov.io/gh/krang-backlink/logger)


```go
func Simple() error {
	opts := logger.NewOptions().
		Service("service").
		Prefix("prefix").
		DefaultStatus("status")

	err := logger.New(context.Background(), opts)
	if err != nil {
		return err
	}

	logger.Info("Hello from Logger!")

	return nil
}
```

```go
func WithWorkplace() error {
	opts := logger.NewOptions().
		Service("api").
		WithWorkplaceNotifier("token", "thread")

	err := logger.New(context.Background(), opts)
	if err != nil {
		return err
	}

	logger.Info("Hello from Logger!")

	return nil
}
```

```go
func WithMongo() error {
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_CONNECTION")).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	opts := logger.NewOptions().
		Service("api").
		WithMongoCollection(client.Database("logs").Collection("col"))

	err = logger.New(context.Background(), opts)
	if err != nil {
		return err
	}

	logger.Info("Hello from Logger!")

	return nil
}
```

```go
func KitchenSink() error {
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_CONNECTION")).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	opts := logger.NewOptions().
		Service("service").
		Prefix("prefix").
		DefaultStatus("status").
		WithWorkplaceNotifier("token", "thread").
		WithMongoCollection(client.Database("logs").Collection("col"))

	err = logger.New(context.Background(), opts)
	if err != nil {
		return err
	}

	logger.Info("Hello from Logger!")

	return nil
}
```
