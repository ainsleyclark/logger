# Logger

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Test](https://github.com/krang-backlink/logger/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/krang-backlink/logger/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/krang-backlink/logger/branch/master/graph/badge.svg?token=HXEZIJ2NUY)](https://codecov.io/gh/krang-backlink/logger)


## Why?

Detailed and verbose logging is important to any application or API. This package aims to make it easy for APIs to log errors to a central location, using a Logrus Hook.

## Installation

```bash
go get -u github.com/krang-backlink/logger
```

## Basic Usage


### Quick Start

![Logger Entries](res/entries.png)

```go
func QuickStart() error {
	err := logger.New(context.TODO(), logger.NewOptions().Service("service"))
	if err != nil {
		return err
	}

	logger.Trace("Trace Entry")
	logger.Debug("Debug Entry")
	logger.Info("Info Entry")
	logger.Warn("Warn Entry")
	logger.Error("Error Entry")
	logger.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error()
	logger.Fatal("Fatal Entry")

	return nil
}
```

### Fields

Fields allow you to log out key value pairs to the logger that will appear under data. The simplest way to use the
logger is simply the package-level exported logger.

```go
func Fields() {
	logger.WithFields(logger.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}
```

## Errors

This package is designed to work with [github.com/ainsleyclark/errors][https://github.com/ainsleyclark/errors] as such
the `WithError` function can be used to log deatiled and rich error data.

```go
func WithError() {
	logger.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error()
}
```

## Middleware

Middleware is provided out of the box in the form of a fire hook. Upon receiving a request from the API,
calling `logger.Fire`will send the log entry to stdout with detailed request information and meta.

```go
func Middleware(r *http.Request) {
	logger.Fire(logger.FireHook{
		Request:      r,
		Status:       http.StatusOK,
		Message:      "Message from API",
		Data:         map[string]any{},
		RequestTime:  time.Now(),
		ResponseTime: time.Now(),
		Latency:      100,
	})
}
```

## Recipes

### Simple
Creates a simple logger with stdout.

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

### WithWorkplace

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

### WithMongo

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

### KitchenSink

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
