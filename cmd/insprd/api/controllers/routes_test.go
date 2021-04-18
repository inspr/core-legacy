package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory/fake"
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
		MemoryManager: fake.MockMemoryManager(nil),
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
			name: "apps",
			want: [...]int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "channels",
			want: [...]int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "channeltypes",
			want: [...]int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusMethodNotAllowed,
			},
		},
		{
			name: "alias",
			want: [...]int{
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
				http.StatusInternalServerError,
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
			ts := httptest.NewServer(testServer.Mux)
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
