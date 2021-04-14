package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_initRoutes(t *testing.T) {
	testServer := &Server{
		mux: http.NewServeMux(),
		rdb: &redisClient,
	}
	testServer.initRoutes()
	defaultMethods := [...]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}
	tests := []struct {
		name string
		want [len(defaultMethods)]int
	}{
		{
			name: "newuser",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "deleteuser",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "updatepwd",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "login",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "refreshtoken",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "wrong_route",
			want: [...]int{
				http.StatusNotFound,
				http.StatusNotFound,
				http.StatusNotFound,
				http.StatusNotFound,
				http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(testServer.mux)
			defer ts.Close()
			client := ts.Client()
			for i, statusCodeResult := range tt.want {
				reqURL := ts.URL + "/" + tt.name

				req, err := http.NewRequest(defaultMethods[i], reqURL, nil)
				if err != nil {
					t.Error("error creating request")
				}

				res, _ := client.Do(req)
				if res.StatusCode != statusCodeResult {
					t.Errorf("Method %v in url %v => got %v, wanted %v",
						defaultMethods[i],
						reqURL,
						res.StatusCode,
						statusCodeResult,
					)
				}
			}
		})
	}
}
