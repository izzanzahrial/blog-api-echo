package elastic

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type Elastic struct {
	client *elasticsearch.Client
	index  string
	alias  string
}

func NewElastic(username, password string, addresses ...string) *Elastic {
	if len(addresses)%2 == 0 {
		log.Fatalf("don't use even number for creating elasticsearch node, you create : %d", len(addresses))
	}

	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("error creating client : %s", err)
	}

	return &Elastic{
		client: es,
	}
}

func (e *Elastic) CreateIndex(index string) error {
	e.index = index
	e.alias = index + "_alias"

	res, err := e.client.Indices.Exists([]string{e.index})
	if err != nil {
		return fmt.Errorf("cannot check index existense: %w", err)
	}

	if res.StatusCode != 404 {
		return fmt.Errorf("error index existence response: %s", res.String())
	}

	res, err = e.client.Indices.Create(e.index)
	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error index creation response: %s", res.String())
	}

	res, err = e.client.Indices.PutAlias([]string{e.index}, e.alias)
	if err != nil {
		return fmt.Errorf("cannot create index alias: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error index alias creation response: %s", res.String())
	}

	return nil
}
