package main_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"product/cache"
	"product/core"
	"product/jwtauth"
	"product/logger"
	"product/repo/mock"
	"product/route"
)

func ginServer() *gin.Engine {

	log := logger.New()

	service := core.NewService(
		mock.NewMock(),
		cache.Redis("localhost"),
		jwtauth.New(),
		nil,
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
			"Should fetch all products with status 200",
			"%s/api/products",
			"GET",
			200,
			nil,
		},
		{
			"Should fetch one product id 1 with status 200",
			"%s/api/products/1",
			"GET",
			200,
			nil,
		},
		{
			"Should add product in mock with status 201",
			"%s/api/products",
			"POST",
			201,
			bytes.NewBuffer([]byte(`{"sku": "olk-yok", "name": "New Product", "price": 12.3456}`)),
		},
		{
			"Should edit product in mock with status 200",
			"%s/api/products",
			"PUT",
			200,
			bytes.NewBuffer([]byte(`{"id": 1, "name": "Edit Product"}`)),
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

			} else if resp.StatusCode != 200 {
				t.Errorf("Expected status code 200, got %v", resp.StatusCode)
			}
			defer resp.Body.Close()
		})
	}
}
