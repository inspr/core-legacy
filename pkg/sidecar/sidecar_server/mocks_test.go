package sidecarserv

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

func Test_mockServer(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockServer(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mockServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mockReader_ReadMessage(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		channel string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &mockReader{
				err: tt.fields.err,
			}
			got, err := mr.ReadMessage(tt.args.channel)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockReader.ReadMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mockReader.ReadMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mockReader_CommitMessage(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		channel string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &mockReader{
				err: tt.fields.err,
			}
			if err := mr.CommitMessage(tt.args.channel); (err != nil) != tt.wantErr {
				t.Errorf("mockReader.CommitMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mockWriter_WriteMessage(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		channel string
		msg     models.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := &mockWriter{
				err: tt.fields.err,
			}
			if err := mw.WriteMessage(tt.args.channel, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("mockWriter.WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
