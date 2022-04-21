package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/izzanzahrial/blog-api-echo/pkg/posting"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreate(t *testing.T) {
	mockService := new(posting.MockService)
	e := echo.New()

	subtest := []struct {
		name         string
		post         repository.Post
		ctx          context.Context
		expectedData webResponse
		expectedErr  int
	}{
		{
			name: "Succesful Handler Create Post",
			post: repository.Post{
				ID:      1,
				Title:   "Test title",
				Content: "Test content",
			},
			ctx: context.Background(),
			expectedData: webResponse{
				Code:   http.StatusCreated,
				Status: "",
				Data: repository.Post{
					ID:      1,
					Title:   "Test title",
					Content: "Test content",
				},
			},
			expectedErr: http.StatusOK,
		},
		{
			name: "Failed Handler Create Post",
			post: repository.Post{
				ID:      2,
				Title:   "Test title",
				Content: "Test content",
			},
			ctx: context.Background(),
			expectedData: webResponse{
				Code:   http.StatusCreated,
				Status: "",
				Data: repository.Post{
					ID:      2,
					Title:   "Test title",
					Content: "Test content",
				},
			},
			expectedErr: http.StatusInternalServerError,
		},
	}

	for _, test := range subtest {
		t.Run(test.name, func(t *testing.T) {
			f := make(url.Values)
			f.Set("title", test.post.Title)
			f.Set("content", test.post.Content)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := NewPostHandler(mockService)

			switch test.name {
			case "Succesful Handler Create Post":
				mockService.On("Create", test.ctx, test.post).Return(test.post, nil).Once()
			case "Failed Handler Create Post":
				mockService.On("Create", test.ctx, test.post).Return(repository.Post{}, nil).Once()
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
