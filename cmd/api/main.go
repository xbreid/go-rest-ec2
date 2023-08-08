package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
	"go-rest-ec2/data"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "3000"

var retryCount int64

const esEndpoint = "https://search-mw-opensearch-c3skgf6txnbq2tqynrd2mzwxmm.us-east-1.es.amazonaws.com"

type Config struct {
	SQS       *sqs.Client
	QueueUrl  string
	Documents data.Documents
	DB        *data.Queries
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Println("Started Rest Service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to Postgres!")
	}

	// TODO: handle settings configs
	queueUrl := os.Getenv("SQS_URL")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create SQS service client
	sqsSvc := sqs.NewFromConfig(cfg)

	signer, err := requestsigner.NewSignerWithService(cfg, "es")
	if err != nil {
		log.Panicf("creating request signer: %v", err)
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{esEndpoint},
		Signer:    signer,
	})
	if err != nil {
		log.Panicf("creating opensearch client: %v", err)
	}

	app := Config{
		SQS:       sqsSvc,
		QueueUrl:  queueUrl,
		Documents: data.EsNew(client),
		DB:        data.New(conn),
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

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := OpenDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			retryCount++
		} else {
			log.Println("Connected to Postgres!")
			return conn
		}

		if retryCount > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
