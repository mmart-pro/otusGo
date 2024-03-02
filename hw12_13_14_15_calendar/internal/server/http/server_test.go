package internalhttp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/app"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/config"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/suite"
)

type WebApiTestSuite struct {
	suite.Suite
	cancel  context.CancelFunc
	baseUrl string
}

func TestWebApiTestSuite(t *testing.T) {
	suite.Run(t, new(WebApiTestSuite))
}

func (suite *WebApiTestSuite) SetupTest() {
	cfg := config.Config{
		LoggerConfig: config.LoggerConfig{
			Level:   "debug",
			LogFile: "",
		},
		HttpConfig: config.EndpointConfig{
			Host: "localhost",
			Port: "8080",
		},
		GrpcConfig: config.EndpointConfig{
			Host: "localhost",
			Port: "5000",
		},
		StorageConfig: config.StorageConfig{
			UseDb: false,
		},
	}
	var ctx context.Context
	ctx, suite.cancel = context.WithCancel(context.Background())
	go func() {
		err := app.NewCalendar(cfg).
			Startup(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
	suite.baseUrl = "http://" + cfg.HttpConfig.GetEndpoint()
}

func (suite *WebApiTestSuite) TestCRUDs() {
	tests := []struct {
		method     string
		url        string
		statusCode int
		body       string
		resp       string
	}{
		{
			method:     "GET",
			url:        "/v1/event/1",
			statusCode: http.StatusInternalServerError,
		},
		{
			method:     "POST",
			url:        "/v1/event",
			statusCode: http.StatusOK,
			body: `{
				"title": "событие тест 1",
				"startDatetime": "2024-02-23T12:00:00Z",
				"endDatetime": "2024-02-25T12:00:00Z",
				"description": "тестовое описание",
				"userId": -1,
				"notifyBeforeMin": 15
			}`,
		},
		{
			method:     "POST",
			url:        "/v1/event",
			statusCode: http.StatusOK,
			body: `{
				"title": "событие тест 2",
				"startDatetime": "2024-01-23T08:00:00Z",
				"endDatetime": "2024-01-25T10:00:00Z",
				"description": "тестовое описание 2",
				"userId": 2,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method:     "POST",
			url:        "/v1/event",
			statusCode: http.StatusInternalServerError,
			body: `{
				"title": "событие FAIL",
				"startDatetime": "2024-02-23T09:00:00Z",
				"endDatetime": "2024-02-25T11:00:00Z",
				"description": "FAIL",
				"userId": 2,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method:     "PUT",
			url:        "/v1/event",
			statusCode: http.StatusOK,
			body: `{
				"id": 1,
				"title": "событие перенесено на март",
				"startDatetime": "2024-03-23T12:00:00Z",
				"endDatetime": "2024-03-25T12:00:00Z",
				"description": "супер акция в марте",
				"userId": 1,
				"notifyBeforeMin": 5
			}`,
		},
		{
			method:     "GET",
			url:        "/v1/event/1",
			statusCode: http.StatusOK,
			resp: `{
				"id": 1,
				"title": "событие перенесено на март",
				"startDatetime": "2024-03-23T12:00:00Z",
				"endDatetime": "2024-03-25T12:00:00Z",
				"description": "супер акция в марте",
				"userId": 1,
				"notifyBeforeMin": 5
			}`,
		},
		{
			method:     "GET",
			url:        "/v1/event/2",
			statusCode: http.StatusOK,
			resp: `{
				"id": 2,
				"title": "событие тест 2",
				"startDatetime": "2024-01-23T08:00:00Z",
				"endDatetime": "2024-01-25T10:00:00Z",
				"description": "тестовое описание 2",
				"userId": 2,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method:     "GET",
			url:        "/v1/event/3",
			statusCode: http.StatusInternalServerError,
		},
		{
			method:     "DELETE",
			url:        "/v1/event/2",
			statusCode: http.StatusOK,
		},
		{
			method:     "DELETE",
			url:        "/v1/event/22",
			statusCode: http.StatusInternalServerError,
		},
		{
			method:     "GET",
			url:        "/v1/event/2",
			statusCode: http.StatusInternalServerError,
		},
	}

	waitForStart(suite.baseUrl + "/health")

	for i, t := range tests {
		suite.Run(fmt.Sprintf("test %d %s %s", i, t.method, strings.ReplaceAll(t.url, "/", ">")), func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			body := bytes.NewBuffer([]byte(t.body))
			req, err := http.NewRequestWithContext(ctx, t.method, suite.baseUrl+t.url, body)

			// запрос прошёл
			suite.Suite.Require().NoError(err)

			client := &http.Client{}
			resp, err := client.Do(req)
			suite.Suite.Require().NoError(err)
			defer resp.Body.Close()

			// статус код ожидаемый
			suite.Suite.Require().Equal(t.statusCode, resp.StatusCode)

			// если ожидаем ответ, то он верный
			if t.resp != "" {
				body, _ := io.ReadAll(io.Reader(resp.Body))

				var actual, expected model.Event

				json.Unmarshal(body, &actual)
				json.Unmarshal([]byte(t.resp), &expected)

				suite.Suite.Require().Equal(expected, actual)
			}
		})
	}

	suite.cancel()
}

func (suite *WebApiTestSuite) TestGetEvents() {
	tests := []struct {
		method      string
		url         string
		body        string
		expectedLen int
	}{
		// setup
		{
			method: "POST",
			url:    "/v1/event",
			body: `{
				"title": "событие 1",
				"startDatetime": "2024-02-21T12:00:00Z",
				"endDatetime": "2024-02-21T13:00:00Z",
				"description": "описание",
				"userId": -1,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method: "POST",
			url:    "/v1/event",
			body: `{
				"title": "событие 2",
				"startDatetime": "2024-02-23T10:00:00Z",
				"endDatetime": "2024-02-23T12:00:00Z",
				"description": "описание",
				"userId": -1,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method: "POST",
			url:    "/v1/event",
			body: `{
				"title": "событие 3",
				"startDatetime": "2024-02-23T12:00:00Z",
				"endDatetime": "2024-02-23T13:00:00Z",
				"description": "описание",
				"userId": -1,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method: "POST",
			url:    "/v1/event",
			body: `{
				"title": "событие 4",
				"startDatetime": "2024-02-14T12:00:00Z",
				"endDatetime": "2024-02-14T13:00:00Z",
				"description": "описание",
				"userId": -1,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method: "POST",
			url:    "/v1/event",
			body: `{
				"title": "событие 5",
				"startDatetime": "2024-02-26T12:00:00Z",
				"endDatetime": "2024-02-26T13:00:00Z",
				"description": "описание",
				"userId": -1,
				"notifyBeforeMin": 0
			}`,
		},
		{
			method: "POST",
			url:    "/v1/event",
			body: `{
				"title": "событие 6",
				"startDatetime": "2024-03-01T12:00:00Z",
				"endDatetime": "2024-03-01T13:00:00Z",
				"description": "описание",
				"userId": -1,
				"notifyBeforeMin": 0
			}`,
		},
		// tests
		{
			method:      "GET",
			url:         "/v1/events/day/2024-02-23",
			expectedLen: 2,
		},
		{
			method:      "GET",
			url:         "/v1/events/week/2024-02-19",
			expectedLen: 3,
		},
		{
			method:      "GET",
			url:         "/v1/events/month/2024-02-01",
			expectedLen: 5,
		},
	}

	waitForStart(suite.baseUrl + "/health")

	for i, t := range tests {
		suite.Run(fmt.Sprintf("test %d %s %s", i, t.method, strings.ReplaceAll(t.url, "/", ">")), func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			body := bytes.NewBuffer([]byte(t.body))
			req, err := http.NewRequestWithContext(ctx, t.method, suite.baseUrl+t.url, body)

			// запрос прошёл
			suite.Suite.Require().NoError(err)

			client := &http.Client{}
			resp, err := client.Do(req)
			suite.Suite.Require().NoError(err)
			defer resp.Body.Close()

			// статус код ожидаемый
			suite.Suite.Require().Equal(http.StatusOK, resp.StatusCode)

			// проверка длины
			if t.expectedLen != 0 {
				body, _ := io.ReadAll(io.Reader(resp.Body))

				// var actual []model.Event
				var actual struct {
					Events []model.Event
				}
				json.Unmarshal(body, &actual)

				suite.Suite.Require().Equal(t.expectedLen, len(actual.Events))
			}
		})
	}

	suite.cancel()
}

// ждём старта web-сервера
func waitForStart(url string) error {
	deadline := time.Now().Add(10 * time.Second)
	for {
		time.Sleep(time.Millisecond * 250)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			return nil // Успешное соединение
		}
		if time.Now().After(deadline) {
			break
		}
	}
	return fmt.Errorf("server took too long to start")
}
