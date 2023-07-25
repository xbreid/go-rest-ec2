package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"net/http"
)

func (app *Config) HelloWorld(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Root endpoint hit"),
		Data:    "Hello World",
	}

	sMInput := &sqs.SendMessageInput{
		MessageBody: aws.String("Root endpoint hit"),
		QueueUrl:    &app.sqsQueue,
	}

	app.SendSQSMessage(ctx, sMInput)

	app.WriteJson(w, http.StatusAccepted, payload)
}
