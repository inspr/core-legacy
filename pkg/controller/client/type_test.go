package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
)

func TestTypeClient_Delete(t *testing.T) {

	type args struct {
		ctx     context.Context
		context string
		name    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete channel test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
			},
			wantErr: false,
		},
		{
			name: "delete channel with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.NewError().BadRequest().Build())
					return
				}

				if r.URL.Path != "/types" {
					t.Errorf("path is not types")
				}

				if r.Method != http.MethodDelete {
					t.Errorf("method is not DELETE")
				}

				var di models.TypeQueryDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}
				if di.TypeName != tt.args.name {
					t.Errorf("name set incorrectly. want = %v, got = %v", di.TypeName, tt.args.name)
				}

				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &TypeClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Delete(tt.args.ctx, tt.args.context, tt.args.name, false); (err != nil) != tt.wantErr {
				t.Errorf("TypeClient.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTypeClient_Get(t *testing.T) {

	type args struct {
		ctx     context.Context
		context string
		name    string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.Type
		wantErr bool
	}{
		{
			name: "get channel test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
			},
			wantErr: false,
			want: &meta.Type{
				Meta: meta.Metadata{
					Name:      "app2",
					Reference: "app1",
				},
				Schema: "schemma",
			},
		},
		{
			name: "get channel with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.NewError().BadRequest().Build())
					return
				}

				if r.URL.Path != "/types" {
					t.Errorf("path is not types")
				}

				if r.Method != http.MethodGet {
					t.Errorf("method is not GET")
				}

				var di models.TypeQueryDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}
				if di.TypeName != tt.args.name {
					t.Errorf("name set incorrectly. want = %v, got = %v", di.TypeName, tt.args.name)
				}

				encoder.Encode(tt.want)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &TypeClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			got, err := ac.Get(tt.args.ctx, tt.args.context, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeClient.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypeClient.Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeClient_Create(t *testing.T) {

	type args struct {
		ctx     context.Context
		context string
		ch      *meta.Type
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test create channel",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch: &meta.Type{
					Meta: meta.Metadata{
						Name: "app3",
					},
					Schema: "schema",
				},
			},
			wantErr: false,
		},
		{
			name: "create channel with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch:      &meta.Type{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.NewError().BadRequest().Build())
					return
				}

				if r.URL.Path != "/types" {
					t.Errorf("path is not types")
				}

				if r.Method != http.MethodPost {
					t.Errorf("method is not POST")
				}

				var di models.TypeDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}

				if !reflect.DeepEqual(di.Type, *tt.args.ch) {
					t.Errorf("request is different. want = \n%+v, \ngot = \n%+v", di.Type, tt.args.ch)
				}
				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &TypeClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Create(tt.args.ctx, tt.args.context, tt.args.ch, false); (err != nil) != tt.wantErr {
				t.Errorf("TypeClient.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTypeClient_Update(t *testing.T) {

	type args struct {
		ctx     context.Context
		context string
		ch      *meta.Type
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test create channel",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch: &meta.Type{
					Meta: meta.Metadata{
						Name: "app3",
					},
					Schema: "schema",
				},
			},
			wantErr: false,
		},
		{
			name: "create channel with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch:      &meta.Type{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.NewError().BadRequest().Build())
					return
				}

				if r.URL.Path != "/types" {
					t.Errorf("path is not types")
				}

				if r.Method != http.MethodPut {
					t.Errorf("method is not PUT")
				}

				var di models.TypeDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}

				if !reflect.DeepEqual(di.Type, *tt.args.ch) {
					t.Errorf("request is different. want = \n%+v, \ngot = \n%+v", di.Type, tt.args.ch)
				}
				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &TypeClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Update(tt.args.ctx, tt.args.context, tt.args.ch, false); (err != nil) != tt.wantErr {
				t.Errorf("TypeClient.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
