package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/izaakdale/lib/listener"
	"github.com/kelseyhightower/envconfig"
)

var (
	name = "service-ticket"
	spec specification
)

type specification struct {
	AWSRegion   string `envconfig:"AWS_REGION"`
	QueueURL    string `envconfig:"QUEUE_URL"`
	AWSEndpoint string `envconfig:"AWS_ENDPOINT"`
}

func main() {
	err := envconfig.Process("", &spec)
	if err != nil {
		panic(err)
	}

	log.Printf("running %s\n", name)
	cfg, err := config.LoadDefaultConfig(context.Background(), func(o *config.LoadOptions) error {
		o.Region = spec.AWSRegion
		return nil
	})

	err = listener.Initialise(cfg, spec.QueueURL, listener.WithEndpoint(spec.AWSEndpoint))
	if err != nil {
		panic(err)
	}

	errChan := make(chan error, 0)
	listener.Listen(func(m listener.Message) error { return nil }, errChan)
}
