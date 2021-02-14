package main_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	// "google.golang.org/grpc"

	"cart/core"
	// "cart/catalogpb"
	"cart/jwtauth"
	"cart/logger"
	"cart/repo/mock"
	"cart/route"
)

func ginServer() *gin.Engine {

	log := logger.New()

	// conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer conn.Close()

	service := core.NewService(
		mock.NewMock(),
		jwtauth.New(),
		// catalogpb.NewCatalogClient(conn),
		nil,
		nil,
	)

	r := route.New(log, true)

	r.Handle(service)

	return r.Engine
}

func TestCart(t *testing.T) {

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
			"Should create a new cart with status 201",
			"%s/api/cart",
			"POST",
			201,
			nil,
		},
		// {
		// 	"Should add an item to the cart with status 200",
		// 	"%s/api/cart/1",
		// 	"POST",
		// 	200,
		// 	bytes.NewBuffer([]byte(`{"id": 1, "qty": 5}`)),
		// },
		{
			"Should fetch all items from cart by id with status 200",
			"%s/api/cart/1",
			"GET",
			200,
			nil,
		},
		{
			"Should update an item in the cart with status 200",
			"%s/api/cart/1",
			"PUT",
			200,
			bytes.NewBuffer([]byte(`{"id": 1, "qty": 1}`)),
		},
		{
			"Should delete product with status 200",
			"%s/api/cart/1/1",
			"DELETE",
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
