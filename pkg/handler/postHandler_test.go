package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/izzanzahrial/blog-api-echo/pkg/posting"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreate(t *testing.T) {
	mockService := new(posting.MockService)
	e := echo.New()

	subtest := []struct {
		status       bool
		name         string
		post         repository.PostData
		ctx          context.Context
		expectedData webResponse
		expectedErr  error
	}{
		{
			status: true,
			name:   "Succesful Handler Create Post",
			post: repository.PostData{
				ID:        1,
				Title:     "Test title",
				ShortDesc: "Test short description",
				Content:   "Test content",
				CreatedAt: time.Now(),
			},
			ctx: context.Background(),
			expectedData: webResponse{
				Code:    http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
				Data: repository.PostData{
					ID:        1,
					Title:     "Test title",
					ShortDesc: "Test short description",
					Content:   "Test content",
					CreatedAt: time.Now(),
				},
			},
			expectedErr: nil,
		},
		{
			status: false,
			name:   "Failed Handler Create Post",
			post: repository.PostData{
				ID:        2,
				Title:     "Test title",
				ShortDesc: "Test short description",
				Content:   "Test content",
			},
			ctx: context.Background(),
			expectedData: webResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusAccepted),
				Data:    echo.ErrInternalServerError,
			},
			expectedErr: echo.ErrInternalServerError,
		},
	}

	for _, test := range subtest {
		t.Run(test.name, func(t *testing.T) {
			f := make(url.Values)
			f.Set("title", test.post.Title)
			f.Set("short_desc", test.post.ShortDesc)
			f.Set("content", test.post.Content)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := NewPostHandler(mockService)

			switch test.status {
			case true:
				mockService.On("Create", test.ctx, test.post).Return(test.post, nil).Once()
			case false:
				mockService.On("Create", test.ctx, test.post).Return(repository.PostData{}, repository.ErrFailedToCreatePost).Once()
			}

			h.Create(c)

			var data webResponse
			dataByte := rec.Body.Bytes()
			err := json.Unmarshal(dataByte, &data)
			if err != nil {
				t.Errorf("Failed to unmarshal data to webresponse")
			}

			assert.Equal(t, test.expectedData, data)
			assert.Equal(t, test.expectedErr, rec.Code)

			mockService.AssertExpectations(t)
		})
	}
}
