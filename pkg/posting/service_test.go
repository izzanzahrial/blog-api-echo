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
	mockDB := new(MockDBtx)
	mockRedis := new(redisDB.MockRedis)
	mockElastic := new(elastic.MockElastic)
	validator := validator.New()

	service := NewService(mockRepo, mockDB, validator, mockRedis, mockElastic)

	subtests := []struct {
		name         string
		ctx          context.Context
		post         PostData
		expectedData PostData
		expectedErr  error
	}{
		{
			name: "Succesful Service Create Post",
			ctx:  context.Background(),
			post: PostData{
				Title:     "Test title",
				ShortDesc: "Test description",
				Content:   "Test content",
			},
			expectedData: PostData{
				Title:     "Test title",
				ShortDesc: "Test description",
				Content:   "Test content",
			},
			expectedErr: nil,
		},
		{
			name: "Failed Service Create Post",
			ctx:  context.Background(),
			post: PostData{
				Title:     "Test title Error",
				ShortDesc: "Test description Error",
				Content:   "Test content Error",
			},
			expectedData: PostData{},
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
					repository.PostData{}, repository.ErrFailedToCreatePost).Once()
			}

			// check return error in mock repo
			// create mock redis on

			data, err := service.Create(test.ctx, test.post)
			assert.Equal(t, test.expectedData, data)
			assert.Equal(t, test.expectedErr, err)

			mockDB.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
