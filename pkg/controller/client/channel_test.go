package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/insprd/api/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"github.com/inspr/inspr/pkg/rest/request"
)

func TestChannelClient_Delete(t *testing.T) {

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

				if r.URL.Path != "/channels" {
					t.Errorf("path is not channels")
				}

				if r.Method != "DELETE" {
					t.Errorf("method is not DELETE")
				}

				var di models.ChannelQueryDI

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if di.Ctx != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", di.Ctx, tt.args.context)
				}
				if di.ChName != tt.args.name {
					t.Errorf("name set incorrectly. want = %v, got = %v", di.ChName, tt.args.name)
				}

				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &ChannelClient{
				c: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Delete(tt.args.ctx, tt.args.context, tt.args.name, false); (err != nil) != tt.wantErr {
				t.Errorf("ChannelClient.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannelClient_Get(t *testing.T) {

	type args struct {
		ctx     context.Context
		context string
		name    string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.Channel
		wantErr bool
	}{
		{
			name: "get channel test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
			},
			wantErr: false,
			want: &meta.Channel{
				Meta: meta.Metadata{
					Name:      "app2",
					Reference: "app1",
				},
				Spec: meta.ChannelSpec{
					Type: "channeltype1",
				},
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

				if r.URL.Path != "/channels" {
					t.Errorf("path is not channels")
				}

				if r.Method != "GET" {
					t.Errorf("method is not GET")
				}

				var di models.ChannelQueryDI

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if di.Ctx != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", di.Ctx, tt.args.context)
				}
				if di.ChName != tt.args.name {
					t.Errorf("name set incorrectly. want = %v, got = %v", di.ChName, tt.args.name)
				}

				encoder.Encode(tt.want)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &ChannelClient{
				c: request.NewJSONClient(s.URL),
			}
			got, err := ac.Get(tt.args.ctx, tt.args.context, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChannelClient.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelClient.Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelClient_Create(t *testing.T) {

	type args struct {
		ctx     context.Context
		context string
		ch      *meta.Channel
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
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name: "app3",
					},
					Spec: meta.ChannelSpec{
						Type: "channeltype1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create channel with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch:      &meta.Channel{},
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

				if r.URL.Path != "/channels" {
					t.Errorf("path is not channels")
				}

				if r.Method != "POST" {
					t.Errorf("method is not POST")
				}

				var di models.ChannelDI

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if di.Ctx != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", di.Ctx, tt.args.context)
				}

				if !reflect.DeepEqual(di.Channel, *tt.args.ch) {
					t.Errorf("request is different. want = \n%+v, \ngot = \n%+v", di.Channel, tt.args.ch)
				}
				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &ChannelClient{
				c: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Create(tt.args.ctx, tt.args.context, tt.args.ch, false); (err != nil) != tt.wantErr {
				t.Errorf("ChannelClient.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannelClient_Update(t *testing.T) {

	type args struct {
		ctx     context.Context
		context string
		ch      *meta.Channel
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
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name: "app3",
					},
					Spec: meta.ChannelSpec{
						Type: "channeltype1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create channel with error test",
			args: args{
				ctx:     context.Background(),
				context: "app1.app2",
				ch:      &meta.Channel{},
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

				if r.URL.Path != "/channels" {
					t.Errorf("path is not channels")
				}

				if r.Method != "PUT" {
					t.Errorf("method is not PUT")
				}

				var di models.ChannelDI

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&di)
				if err != nil {
					t.Error(err)
				}

				if di.Ctx != tt.args.context {
					t.Errorf("context set incorrectly. want = %v, got = %v", di.Ctx, tt.args.context)
				}

				if !reflect.DeepEqual(di.Channel, *tt.args.ch) {
					t.Errorf("request is different. want = \n%+v, \ngot = \n%+v", di.Channel, tt.args.ch)
				}
				encoder.Encode(diff.Changelog{})
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &ChannelClient{
				c: request.NewJSONClient(s.URL),
			}
			if _, err := ac.Update(tt.args.ctx, tt.args.context, tt.args.ch, false); (err != nil) != tt.wantErr {
				t.Errorf("ChannelClient.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
