package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
)

type Elastic struct {
	Client *elasticsearch.Client
	Index  string
	Alias  string
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
		Client: es,
	}
}

func (e *Elastic) CreateIndex(index string) error {
	e.Index = index
	e.Alias = index + "_alias"

	res, err := e.Client.Indices.Exists([]string{e.Index})
	if err != nil {
		return fmt.Errorf("cannot check index existense: %w", err)
	}

	if res.StatusCode != 404 {
		return fmt.Errorf("error index existence response: %s", res.String())
	}

	res, err = e.Client.Indices.Create(e.Index)
	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error index creation response: %s", res.String())
	}

	res, err = e.Client.Indices.PutAlias([]string{e.Index}, e.Alias)
	if err != nil {
		return fmt.Errorf("cannot create index alias: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error index alias creation response: %s", res.String())
	}

	return nil
}

func (e *Elastic) Insert(ctx context.Context, post repository.Post) error {
	body, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      e.Index,
		DocumentID: strconv.Itoa(int(post.ID)),
		Body:       bytes.NewReader(body),
	}

	res, err := req.Do(ctx, e.Client)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed because there's an error in response: %s", res.String())
	}

	return nil
}

func (e *Elastic) Update(ctx context.Context, post repository.Post) error {
	body, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      e.Index,
		DocumentID: strconv.Itoa(int(post.ID)),
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}

	res, err := req.Do(ctx, e.Client)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed because there's an error in response: %s", res.String())
	}

	return nil
}

func (e *Elastic) Delete(ctx context.Context, postID string) error {
	req := esapi.DeleteRequest{
		Index:      e.Index,
		DocumentID: postID,
	}

	res, err := req.Do(ctx, e.Client)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed because there's an error in response: %s", res.String())
	}

	return nil
}
