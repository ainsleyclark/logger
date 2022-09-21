// Copyright 2022 Ainsley Clark. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package examples

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/logger"
	"github.com/ainsleyclark/logger/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// Simple godoc
// Creates a simple logger with stdout.
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

// WithWorkplace godoc
// Create a logger with Facebook Workplace integration. A token and a
// thread are required to send any error code that has been marked
// as `errors.INTERNAL` to thread ID passed.
func WithWorkplace() error {
	opts := logger.NewOptions().
		Service("api").
		WithWorkplaceNotifier("token", "thread", nil, nil)

	err := logger.New(context.Background(), opts)
	if err != nil {
		return err
	}

	logger.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error()

	return nil
}

// WithWorkplaceReport godoc
// You can pass a function to `WithWorkplaceNotifier` as the second argument
// which is a callback function to determine if the log entry should be
// sent to a thread, an example is below.
func WithWorkplaceReport() {
	// Don't send the message to workplace if there is no error.
	workplaceCallBack := func(entry types.Entry) bool {
		if !entry.HasError() {
			return false
		}
		return true
	}

	_ = logger.NewOptions().
		Service("api").
		WithWorkplaceNotifier("token", "thread", workplaceCallBack, nil)

	// etc
}

// WithWorkplaceFormatter godoc
// You can pass a function to `WithWorkplaceNotifier` as the third argument
// which is a callback function to write the message to Workplace.
// This is where you can customise the message easily and return
// a formatted string.
func WithWorkplaceFormatter() {
	// Format the message with the supplied arguments.
	workplaceFormatter := func(entry types.Entry, args types.FormatMessageArgs) string {
		return args.Version + " " + " hello from Workplace!"
	}

	_ = logger.NewOptions().
		Service("api").
		WithWorkplaceNotifier("token", "thread", nil, workplaceFormatter)

	// etc
}

// WithMongo godoc
// Create a logger with Mongo integration. All logs are sent to the
// collection passed.
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
		WithMongoCollection(client.Database("logs").Collection("col"), nil)

	err = logger.New(context.Background(), opts)
	if err != nil {
		return err
	}

	logger.Info("Hello from Logger!")

	return nil
}

// WithMongoReport godoc
// You can pass a function to `WithMongoCollection` as the second argument
// which is a callback function to determine if the message should be
// stored within Mongo, an example is below.
func WithMongoReport() {
	// Don't send the message to Mongo if there is no error.
	mongoCallBack := func(entry types.Entry) bool {
		if !entry.HasError() {
			return false
		}
		return true
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION")))
	if err != nil {
		log.Fatalln(err)
	}

	_ = logger.NewOptions().
		Service("api").
		WithMongoCollection(client.Database("logs").Collection("col"), mongoCallBack)

	// etc
}

// KitchenSink godoc
// Boostrap all Log integrations.
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
		WithWorkplaceNotifier("token", "thread", nil, nil).
		WithMongoCollection(client.Database("logs").Collection("col"), nil)

	err = logger.New(context.Background(), opts)
	if err != nil {
		return err
	}

	logger.Info("Hello from Logger!")

	return nil
}
