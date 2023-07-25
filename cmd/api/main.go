package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"net/http"
	"os"
)

const webPort = "3000"

type Config struct {
	sqsClient *sqs.Client
	sqsQueue  string
}

func main() {
	log.Println("Started Rest Service")

	// TODO: handle settings configs
	queueUrl := os.Getenv("SQS_URL")
	log.Printf("QUEUE_URL: %s", queueUrl)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create SQS service client
	sqsSvc := sqs.NewFromConfig(cfg)

	app := Config{
		sqsClient: sqsSvc,
		sqsQueue:  queueUrl,
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
