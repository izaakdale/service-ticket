package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/izaakdale/lib/listener"
	"github.com/izaakdale/lib/server"
	"github.com/izaakdale/service-event-order/pkg/proto/order"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name   = "service-ticket"
	spec   specification
	client order.OrderServiceClient
)

type specification struct {
	Host        string `envconfig:"HOST"`
	Port        string `envconfig:"PORT"`
	AWSRegion   string `envconfig:"AWS_REGION"`
	QueueURL    string `envconfig:"QUEUE_URL"`
	AWSEndpoint string `envconfig:"AWS_ENDPOINT"`
	GRPCHost    string `envconfig:"GRPC_HOST"`
	GRPCPort    string `envconfig:"GRPC_PORT"`
}

func Run() {
	err := envconfig.Process("", &spec)
	if err != nil {
		panic(err)
	}

	log.Printf("running %s\n", name)
	cfg, err := config.LoadDefaultConfig(context.Background(), func(o *config.LoadOptions) error {
		o.Region = spec.AWSRegion
		return nil
	})

	grpcSocket := fmt.Sprintf("%s:%s", spec.GRPCHost, spec.GRPCPort)

	conn, err := grpc.Dial(grpcSocket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to %s", grpcSocket)
	}
	defer conn.Close()

	client = order.NewOrderServiceClient(conn)

	err = listener.Initialise(cfg, spec.QueueURL, listener.WithEndpoint(spec.AWSEndpoint))
	if err != nil {
		panic(err)
	}

	srv, err := server.New(
		Router(),
		server.WithHost(spec.Host),
		server.WithPort(spec.Port),
		server.WithTimeouts(time.Second, time.Second))
	if err != nil {
		panic(err)
	}

	go srv.ListenAndServe()

	errChan := make(chan error, 0)

	go listener.Listen(Process, errChan)

	shutCh := make(chan os.Signal, 1)
	signal.Notify(shutCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case err = <-errChan:
		// TODO do we definitely want to quit?
		if err != nil {
			log.Fatal(err)
		}
	case signal := <-shutCh:
		log.Printf("got shutdown signal: %s, exiting\n", signal)
		os.Exit(1)
	}
}
