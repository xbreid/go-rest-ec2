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

	queueName := "sqs-mw-queue"
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	result, err := app.GetQueueURL(ctx, gQInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	fmt.Println(queueURL)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Root endpoint hit"),
		Data:    "Hello World",
	}

	sMInput := &sqs.SendMessageInput{
		MessageBody: aws.String("Information about the NY Times fiction bestseller for the week of 12/11/2016."),
		QueueUrl:    queueURL,
	}

	app.SendSQSMessage(ctx, sMInput)

	app.WriteJson(w, http.StatusAccepted, payload)
}
