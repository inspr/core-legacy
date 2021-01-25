package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/mocks"
)

// TestServer_initRoutes - this test is a bit different than the one automatically
// generated, the idea behind it is to specify in wanted the desired result for each
// of the 4 default methods [GET,POST,PUT,DELETE] being a 405 a invalid request. It is
// important to make clear that when the proper method is used the desired http response
// is the StatusBadRequest(400) due to not putting values in the body of
// the requests
func TestServer_initRoutes(t *testing.T) {
	testServer := &Server{
		Mux:           http.NewServeMux(),
		MemoryManager: mocks.MockMemoryManager(nil),
	}
	testServer.initRoutes()
	defaultMethods := [...]string{"GET", "POST", "PUT", "DELETE"}
	tests := []struct {
		name string
		want [len(defaultMethods)]int
	}{
		{
			name: "apps",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusBadRequest,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "apps/ref",
			want: [...]int{
				http.StatusBadRequest,
				http.StatusMethodNotAllowed,
				http.StatusBadRequest,
				http.StatusBadRequest,
			},
		},
		{
			name: "channels",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusBadRequest,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "channels/ref",
			want: [...]int{
				http.StatusBadRequest,
				http.StatusMethodNotAllowed,
				http.StatusBadRequest,
				http.StatusBadRequest,
			},
		},
		{
			name: "channeltypes",
			want: [...]int{
				http.StatusMethodNotAllowed,
				http.StatusBadRequest,
				http.StatusMethodNotAllowed,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "channeltypes/ref",
			want: [...]int{
				http.StatusBadRequest,
				http.StatusMethodNotAllowed,
				http.StatusBadRequest,
				http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(testServer.Mux)
			defer ts.Close()
			client := ts.Client()
			for i, statusCodeResult := range tt.want {
				reqURL := ts.URL + "/" + tt.name
				var req *http.Request
				var err error
				switch defaultMethods[i] {
				case "GET":
					req, err = http.NewRequest(http.MethodGet, reqURL, nil)
				case "POST":
					req, err = http.NewRequest(http.MethodPost, reqURL, nil)
				case "PUT":
					req, err = http.NewRequest(http.MethodPut, reqURL, nil)
				case "DELETE":
					req, err = http.NewRequest(http.MethodDelete, reqURL, nil)
				default:
				}
				if err != nil {
					t.Error("error creating request")
				}
				res, err := client.Do(req)
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
