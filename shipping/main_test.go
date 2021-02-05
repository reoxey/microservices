package main_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"shipping/jwtauth"
	"shipping/logger"
	"shipping/repo/mock"
	"shipping/route"
	"shipping/core"
)

func ginServer() *gin.Engine {

	log := logger.New()

	service := core.NewService(
		mock.NewMock(),
		jwtauth.New(),
	)

	r := route.New(log, true)

	r.Handle(service)

	return r.Engine
}

func TestProduct(t *testing.T) {

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
			"Should fetch all addresses with status 200",
			"%s/api/shipping/addresses",
			"GET",
			200,
			nil,
		},
		{
			"Should fetch one address id 1 with status 200",
			"%s/api/shipping/addresses/1",
			"GET",
			200,
			nil,
		},
		{
			"Should add address in mock with status 201",
			"%s/api/shipping/addresses",
			"POST",
			201,
			bytes.NewBuffer([]byte(`{"user": 1, "contact_name": "Test", "contact_phone": "9456", "landmark": "eins", "city": "zwei", "state": "drei", "country": "vier", "zip": 123456}`)),
		},
		{
			"Should edit address in mock with status 200",
			"%s/api/shipping/addresses",
			"PUT",
			200,
			bytes.NewBuffer([]byte(`{"id": 1, "contact_name": "Test zwei", "contact_phone": "1234566689"}`)),
		},
		{
			"Should update shipping status with 200 code",
			"%s/api/shipping/status/1/1",
			"PUT",
			200,
			nil,
		},
		{
			"Should update payment status with 200 code",
			"%s/api/shipping/payment/1/1",
			"PUT",
			200,
			nil,
		},
	}

	token, err := jwtauth.New().GenerateToken(1, "", true)
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
				t.Errorf("Expected status code 200, got %v", resp.StatusCode)
			}
			defer resp.Body.Close()
		})
	}
}
