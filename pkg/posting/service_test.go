package posting

import (
	"context"
	"testing"

	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {
	service := service{
		Repository: new(repository.MockPostgre),
		DB:         new(MockDB),
		Validate:   new(MockValidator),
	}
	subtests := []struct {
		name         string
		ctx          context.Context
		post         repository.Post
		expectedData repository.Post
		expectedErr  error
	}{
		{
			name: "Succesful Create Post",
			ctx:  context.Background(),
			post: repository.Post{
				ID:      1,
				Title:   "Test title",
				Content: "Test content",
			},
			expectedData: repository.Post{
				ID:      1,
				Title:   "Test title",
				Content: "Test content",
			},
			expectedErr: nil,
		},
		{
			name: "Failed Create Post",
			ctx:  context.Background(),
			post: repository.Post{
				ID:      1,
				Title:   "Test title",
				Content: "Test content",
			},
			expectedData: repository.Post{
				ID:      1,
				Title:   "Test title",
				Content: "Test content",
			},
			expectedErr: repository.ErrFailedToAddPost,
		},
	}

	for _, test := range subtests {
		t.Run(test.name, func(t *testing.T) {
			data, err := service.Create(test.ctx, test.post)
			assert.Equal(t, test.expectedData, data)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
