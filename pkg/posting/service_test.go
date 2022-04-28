package posting

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/pkg/elastic"
	redisDB "github.com/izzanzahrial/blog-api-echo/pkg/redis"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func TestServiceCreate(t *testing.T) {
	mockRepo := new(repository.MockPostgre)
	mockDB := new(MockDB)
	mockRedis := new(redisDB.MockRedis)
	mockElastic := new(elastic.MockElastic)
	validator := validator.New()

	service := NewService(mockRepo, mockDB, validator, mockRedis, mockElastic)

	subtests := []struct {
		name         string
		ctx          context.Context
		post         repository.Post
		expectedData repository.Post
		expectedErr  error
	}{
		{
			name: "Succesful Service Create Post",
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
			name: "Failed Service Create Post",
			ctx:  context.Background(),
			post: repository.Post{
				ID:      2,
				Title:   "Test title Error",
				Content: "Test content Error",
			},
			expectedData: repository.Post{},
			expectedErr:  repository.ErrFailedToCreatePost,
		},
	}

	for _, test := range subtests {
		t.Run(test.name, func(t *testing.T) {
			mockDB.On("BEGIN").Return(&sql.Tx{}, nil).Once()
			switch test.name {
			case "Succesful Service Create Post":
				mockRepo.On("Create", test.ctx, &sql.Tx{}, test.post).Return(test.post, nil).Once()
			case "Failed Service Create Post":
				mockRepo.On("Create", test.ctx, &sql.Tx{}, test.post).Return(
					repository.Post{}, repository.ErrFailedToCreatePost).Once()
			}

			data, err := service.Create(test.ctx, test.post)
			assert.Equal(t, test.expectedData, data)
			assert.Equal(t, test.expectedErr, err)

			mockDB.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
