package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"net/http"
)

type messagePayload struct {
	Host      string      `json:"host"`
	Operation string      `json:"operation"`
	Resource  string      `json:"resource"`
	Payload   interface{} `json:"payload"`
}

type AccountGroupPayload struct {
	DisplayName   string `json:"display_name"`
	Country       string `json:"country"`
	Locality      string `json:"locality"`
	PostalCode    string `json:"postal_code"`
	StreetAddress string `json:"street_address"`
	Region        string `json:"region"`
	ExternalID    string `json:"external_id"`
	Active        bool   `json:"active"`
}

func (app *Config) PgSearch(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	search := r.URL.Query().Get("search")

	if search == "" {
		app.WriteJson(w, http.StatusBadRequest, "search parameter required")
		return
	}

	accountGroups, err := app.DB.Search(ctx, search)
	if err != nil {
		log.Println(err)
		app.WriteJson(w, http.StatusInternalServerError, "pg search failure")
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "accepted",
		Data:    accountGroups,
	}

	app.WriteJson(w, http.StatusAccepted, payload)
}

func (app *Config) EsSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	if search == "" {
		app.WriteJson(w, http.StatusBadRequest, "search parameter required")
		return
	}

	accountGroups, err := app.Documents.AccountGroupSearch(search)
	if err != nil {
		log.Println(err)
		app.WriteJson(w, http.StatusInternalServerError, "elastic search failure")
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "accepted",
		Data:    accountGroups,
	}

	app.WriteJson(w, http.StatusAccepted, payload)
}

func (app *Config) UpsertAccountGroup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var requestPayload AccountGroupPayload

	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		app.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("requestPayload: %v", requestPayload)

	if requestPayload.ExternalID == "" {
		app.WriteJson(w, http.StatusBadRequest, "external id required for upsert")
		return
	}

	msgPayload := messagePayload{
		Host:      "go-rest-fargate",
		Operation: "UPSERT",
		Resource:  "ACCOUNT_GROUP",
		Payload:   requestPayload,
	}

	body, err := json.Marshal(msgPayload)
	if err != nil {
		app.WriteJson(w, http.StatusInternalServerError, "msg payload failure")
		return
	}

	sMInput := &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    &app.QueueUrl,
	}

	app.SendSQSMessage(ctx, sMInput)

	payload := jsonResponse{
		Error:   false,
		Message: "queued",
		Data:    nil,
	}

	app.WriteJson(w, http.StatusAccepted, payload)
}

func (app *Config) HelloWorld(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hello World",
		Data:    nil,
	}

	app.WriteJson(w, http.StatusAccepted, payload)
}
