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

func TestAliasClient_Get(t *testing.T) {
	type args struct {
		ctx     context.Context
		context string
		name    string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.Alias
		wantErr bool
	}{
		{
			name: "get alias test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
			},
			wantErr: false,
			want: &meta.Alias{
				Meta: meta.Metadata{
					Name:      "app2",
					Reference: "app1",
				},
				Target: "alias_target",
			},
		},
		{
			name: "get alias with error test",
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

				if r.URL.Path != "/alias" {
					t.Errorf("path is not alias")
				}

				if r.Method != http.MethodGet {
					t.Errorf("method is not GET")
				}

				var di models.AliasQueryDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}
				if di.Key != tt.args.name {
					t.Errorf("name set incorrectly. want = %v, got = %v", di.Key, tt.args.name)
				}

				encoder.Encode(tt.want)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &AliasClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			got, err := ac.Get(tt.args.ctx, tt.args.context, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AliasClient.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AliasClient.Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliasClient_Create(t *testing.T) {
	type args struct {
		ctx     context.Context
		context string
		ch      *meta.Alias
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test create alias",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch: &meta.Alias{
					Meta: meta.Metadata{
						Name: "alias_name",
					},
					Target: "alias_target",
				},
			},
			wantErr: false,
		},
		{
			name: "create alias with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch:      &meta.Alias{},
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

				if r.URL.Path != "/alias" {
					t.Errorf("path is not alias")
				}

				if r.Method != http.MethodPost {
					t.Errorf("method is not POST")
				}

				var di models.AliasDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}

				if !reflect.DeepEqual(di.Alias, *tt.args.ch) {
					t.Errorf("request is different. want = \n%+v, \ngot = \n%+v", di.Alias, tt.args.ch)
				}
				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &AliasClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Create(tt.args.ctx, tt.args.context, "alias_target", tt.args.ch, false); (err != nil) != tt.wantErr {
				t.Errorf("AliasClient.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasClient_Delete(t *testing.T) {
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
			name: "delete alias test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
			},
			wantErr: false,
		},
		{
			name: "delete alias with error test",
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

				if r.URL.Path != "/alias" {
					t.Errorf("path is not alias")
				}

				if r.Method != http.MethodDelete {
					t.Errorf("method is not DELETE")
				}

				var di models.AliasQueryDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}
				if di.Key != tt.args.name {
					t.Errorf("name set incorrectly. want = %v, got = %v", di.Key, tt.args.name)
				}

				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &AliasClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Delete(tt.args.ctx, tt.args.context, tt.args.name, false); (err != nil) != tt.wantErr {
				t.Errorf("AliasClient.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasClient_Update(t *testing.T) {
	type args struct {
		ctx     context.Context
		context string
		ch      *meta.Alias
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test update alias",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch: &meta.Alias{
					Meta: meta.Metadata{
						Name: "alias_name",
					},
					Target: "alias_target",
				},
			},
			wantErr: false,
		},
		{
			name: "update alias with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch:      &meta.Alias{},
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

				if r.URL.Path != "/alias" {
					t.Errorf("path is not alias")
				}

				if r.Method != "PUT" {
					t.Errorf("method is not PUT")
				}

				var di models.AliasDI
				scope := r.Header.Get(rest.HeaderScopeKey)

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if scope != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", scope, tt.args.context)
				}

				if !reflect.DeepEqual(di.Alias, *tt.args.ch) {
					t.Errorf("request is different. want = \n%+v, \ngot = \n%+v", di.Alias, tt.args.ch)
				}
				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &AliasClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Update(tt.args.ctx, tt.args.context, "alias_target", tt.args.ch, false); (err != nil) != tt.wantErr {
				t.Errorf("AliasClient.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
