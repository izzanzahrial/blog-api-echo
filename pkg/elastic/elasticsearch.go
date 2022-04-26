package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

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

type document struct {
	Source interface{} `json:"_source"`
}

func (e *Elastic) FindByID(ctx context.Context, postID string) (repository.Post, error) {
	req := esapi.GetRequest{
		Index:      e.Index,
		DocumentID: postID,
	}

	res, err := req.Do(ctx, e.Client)
	if err != nil {
		return repository.Post{}, fmt.Errorf("failed to find the post by id: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return repository.Post{}, fmt.Errorf("failed because there's an error in response: %s", res.String())
	}

	var (
		post repository.Post
		body document
	)
	body.Source = &post

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return repository.Post{}, fmt.Errorf("failed to decode the result body: %w", err)
	}

	return post, nil
}

func (e *Elastic) SearchPost(ctx context.Context, query string, from int, size int) error {

	res, err := e.Client.Search(
		e.Client.Search.WithContext(ctx),
		e.Client.Search.WithIndex(e.Index),
		e.Client.Search.WithBody(e.BuildBody(from, size, query)),
		e.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	// look at the result from elasticsearch
	// make the spesific struct for that result
}

func (e *Elastic) BuildBody(from int, size int, query string) io.Reader {
	var body strings.Builder

	body.WriteString("{\n")

	if query == "" {
		body.WriteString(searchAll)
	} else {
		body.WriteString(fmt.Sprintf(searchMatch, from, size, query))
	}

	body.WriteString("\n}")

	return strings.NewReader(body.String())
}

const searchAll = `"query": { 
							"match_all": {}
					},
					"size": 25,
					"sort": {
							"title": "asc"
					}`

const searchMatch = `"query": {
							"from": %d,
							"size": %d,
							"multi_match": {
										"query": %q,
										"fields": [
												"title^2",
												"content"
										],
										"type": "phrase"
									}
								}`
