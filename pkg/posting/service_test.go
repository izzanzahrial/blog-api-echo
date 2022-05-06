package posting

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/izzanzahrial/blog-api-echo/pkg/elastic"
	redisDB "github.com/izzanzahrial/blog-api-echo/pkg/redis"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func TestServiceCreate(t *testing.T) {
	mockRepo := new(repository.MockPostingPostgre)
	mockDB := new(MockDBtx)
	mockRedis := new(redisDB.MockRedis)
	mockElastic := new(elastic.MockElastic)
	validator := validator.New()

	service := NewService(mockRepo, mockDB, validator, mockRedis, mockElastic)

	subtests := []struct {
		status       bool
		name         string
		ctx          context.Context
		post         PostData
		expectedData PostData
		expectedErr  error
	}{
		{
			status: true,
			name:   "Succesful Service Create Post",
			ctx:    context.Background(),
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
			status: false,
			name:   "Failed Service Create Post",
			ctx:    context.Background(),
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
			mockDB.On("Begin").Return(&sql.Tx{}, nil).Once()
			switch test.status {
			case true:
				mockRepo.On("Create", test.ctx, &sql.Tx{}, test.post).Return(test.post, nil).Once()
			case false:
				mockRepo.On("Create", test.ctx, &sql.Tx{}, test.post).Return(
					repository.PostData{}, repository.ErrFailedToCreatePost).Once()
			}

			// return repo still no created time
			// mock for rollback and commit

			if test.status == true {
				mockRedis.On("Set").Return(&redis.StatusCmd{}).Once()
			}

			data, err := service.Create(test.ctx, test.post)
			assert.Equal(t, test.expectedData, data)
			assert.Equal(t, test.expectedErr, err)

			mockDB.AssertExpectations(t)
			mockRepo.AssertExpectations(t)

			if test.status == true {
				mockRedis.AssertExpectations(t)
			}
		})
	}
}
