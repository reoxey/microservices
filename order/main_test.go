package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"order/core"
	"order/jwtauth"
	"order/logger"
	"order/repo/mock"
	"order/route"
)

func ginServer() *gin.Engine {

	log := logger.New()

	service := core.NewService(
		mock.NewMock(),
		jwtauth.New(),
		nil,
	)

	r := route.New(log, true)

	r.Handle(service)

	return r.Engine
}

func TestOrder(t *testing.T) {

	client := &http.Client{}

	ts := httptest.NewServer(ginServer())

	defer ts.Close()

	tests := []struct {
		reason string
		endpoint string
		method string
		status int
		payload io.Reader
	}{
		{
			"Should fetch an order by id with status 200",
			"%s/api/orders/1",
			"GET",
			200,
			nil,
		},
		{
			"Should fetch all order by buyer id with status 200",
			"%s/api/orders",
			"GET",
			200,
			nil,
		},
	}

	token, err := jwtauth.New().GenerateToken(0, "", true)
	if err != nil {
		t.Error(err)
		return
	}

	for _, st := range tests {

		t.Run(st.reason, func(t *testing.T) {

			req, err := http.NewRequest(st.method, fmt.Sprintf(st.endpoint, ts.URL), st.payload)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)

			} else if resp.StatusCode != st.status {
				t.Errorf("Expected status code %d, got %v", st.status, resp.StatusCode)
			}
			defer resp.Body.Close()
		})
	}
}
