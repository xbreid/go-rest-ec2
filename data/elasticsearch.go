package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"log"
	"strings"
)

var client *opensearch.Client

func EsNew(esClient *opensearch.Client) Documents {
	client = esClient

	return Documents{
		AccountGroup: AccountGroupDoc{},
	}
}

type SearchResponse struct {
	Took int `json:"took"`
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Index  string      `json:"_index"`
			Type   string      `json:"_type"`
			ID     string      `json:"_id"`
			Score  float64     `json:"_score"`
			Source interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Documents struct {
	AccountGroup AccountGroupDoc
}

type AccountGroupDoc struct {
	Uuid          string `json:"uuid"`
	DisplayName   string `json:"display_name"`
	Country       string `json:"country"`
	Locality      string `json:"locality"`
	PostalCode    string `json:"postal_code"`
	StreetAddress string `json:"street_address"`
	Region        string `json:"region"`
	ExternalID    string `json:"external_id"`
	Active        string `json:"active"`
	CreatedAt     string `json:"created_at"`
}

func (d *Documents) AccountGroupSearch(query string) ([]*AccountGroupDoc, error) {
	ctx := context.Background()

	content := fmt.Sprintf(`{
    "size": 25,
		"query": {
			"query_string": {
				"query": "%s"
			}
		}
	}`, query)

	search := opensearchapi.SearchRequest{
		Index: []string{"account_groups"},
		Body:  strings.NewReader(content),
	}

	res, err := search.Do(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error executing search response, %v", err)
	}
	defer res.Body.Close()

	var searchResponse SearchResponse
	if err = json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("error decoding search response %v", err)
	}

	log.Printf("Took: %d ms", searchResponse.Took)
	log.Printf("Total Hits: %d", searchResponse.Hits.Total.Value)

	var accountGroups []*AccountGroupDoc

	for _, hit := range searchResponse.Hits.Hits {
		sourceBytes, err := json.Marshal(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("error marshaling _source %v", err)
		}

		var accountGroup AccountGroupDoc
		if err := json.Unmarshal(sourceBytes, &accountGroup); err != nil {
			return nil, fmt.Errorf("error unmarshaling source doc: %v", err)
		}

		accountGroups = append(accountGroups, &accountGroup)
	}

	return accountGroups, nil
}
