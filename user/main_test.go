package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"user/cache"
	"user/jwtauth"
	"user/logger"
	"user/core"
	"user/repo/mock"
	"user/route"
)

func ginServer() *gin.Engine {

	log := logger.New()

	service := core.NewService(
		mock.NewMock(),
		cache.Redis("localhost"),
		jwtauth.New(),
	)

	r := route.New(log, true)

	r.Handle(service)

	return r.Engine
}

func TestUser(t *testing.T) {

	type auth struct {
		Token string
	}
	login := auth{}

	client := &http.Client{}

	ts := httptest.NewServer(ginServer())

	defer ts.Close()

	tests := []struct {
		reason string
		endpoint string
		method string
		status int
		payload io.Reader
		isLogin	bool
	}{
		{
			"Should login with auth",
			"%s/api/login",
			"POST",
			200,
			bytes.NewBuffer([]byte(`{"email": "", "password": ""}`)),
			true,
		},
		{
			"Should fetch all users with status 200",
			"%s/api/users",
			"GET",
			200,
			nil,
			false,
		},
		{
			"Should fetch one users id 1 with status 200",
			"%s/api/users/1",
			"GET",
			200,
			nil,
			false,
		},
		{
			"Should add user in mock with status 201",
			"%s/api/users",
			"POST",
			201,
			bytes.NewBuffer([]byte(`{"name": "New User", "email": "user@test.com","password": "secret"}`)),
			false,
		},
		{
			"Should edit user in mock with status 200",
			"%s/api/users",
			"PUT",
			200,
			bytes.NewBuffer([]byte(`{"id": 1, "name": "Edit User"}`)),
			false,
		},
	}

	for _, st := range tests {

		t.Run(st.reason, func(t *testing.T) {

			req, err := http.NewRequest(st.method, fmt.Sprintf(st.endpoint, ts.URL), st.payload)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !st.isLogin {
				req.Header.Set("Authorization", "Bearer "+login.Token)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)

			} else if resp.StatusCode != st.status {
				t.Errorf("Expected status code 200, got %v", resp.StatusCode)
			}
			defer resp.Body.Close()

			if st.isLogin {
				err = json.NewDecoder(resp.Body).Decode(&login)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
