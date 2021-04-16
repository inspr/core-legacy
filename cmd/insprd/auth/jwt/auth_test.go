// package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func TestNewJWTauth(t *testing.T) {
	os.Setenv("AUTH_PATH", "mock_url")
	defer func() { os.Remove("AUTH_PATH") }()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	tests := []struct {
		name string
		want *JWTauth
	}{
		{
			name: "returns_JWT_auth",
			want: &JWTauth{
				PublicKey: &privKey.PublicKey,
				authURL:   "mock_url",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewJWTauth(&privKey.PublicKey)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJWTauth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJWTauth_Validate(t *testing.T) {
	os.Setenv("AUTH_PATH", "mock_url")
	defer func() { os.Remove("AUTH_PATH") }()

	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	const aLongLongTimeAgo = 233431200

	invalidToken := func() []byte {
		token := jwt.New()
		token.Set(`foo`, `bar`)
		signed, _ := jwt.Sign(token, jwa.RS256, "differentkey")
		return signed
	}
	expiredToken := func() []byte {
		token := jwt.New()
		token.Set(jwt.ExpirationKey, time.Unix(aLongLongTimeAgo, 0))
		signed, _ := jwt.Sign(token, jwa.RS256, privKey)
		return signed
	}
	nilExpiredToken := func() []byte {
		token := jwt.New()
		token.Set(jwt.ExpirationKey, nil)
		signed, _ := jwt.Sign(token, jwa.RS256, privKey)
		return signed
	}
	noPayloadToken := func() []byte {
		token := jwt.New()
		token.Set(jwt.ExpirationKey, time.Now().Add(30*time.Minute))
		signed, _ := jwt.Sign(token, jwa.RS256, privKey)
		return signed
	}
	fineToken := func() []byte {
		token := jwt.New()
		token.Set(jwt.ExpirationKey, time.Now().Add(30*time.Minute))

		payload := models.Payload{
			UID:        "mock_UID",
			Role:       0,
			Scope:      []string{"mock"},
			Refresh:    []byte("mock_refresh"),
			RefreshURL: "mock_refresh_url",
		}
		token.Set("payload", payload)
		signed, _ := jwt.Sign(token, jwa.RS256, privKey)
		return signed
	}

	type args struct {
		token []byte
	}
	tests := []struct {
		name    string
		JA      *JWTauth
		args    args
		want    *models.Payload
		want1   []byte
		wantErr bool
	}{
		{
			name: "Invalid_token",
			JA:   NewJWTauth(&privKey.PublicKey),
			args: args{
				token: invalidToken(),
			},
			want:    nil,
			want1:   invalidToken(),
			wantErr: true,
		},
		{
			name: "Expired_token",
			JA:   NewJWTauth(&privKey.PublicKey),
			args: args{
				token: expiredToken(),
			},
			want:    nil,
			want1:   expiredToken(),
			wantErr: true,
		},
		{
			name: "nil_Expired_token",
			JA:   NewJWTauth(&privKey.PublicKey),
			args: args{
				token: nilExpiredToken(),
			},
			want:    nil,
			want1:   nilExpiredToken(),
			wantErr: true,
		},
		{
			name: "Payload_notFound",
			JA:   NewJWTauth(&privKey.PublicKey),
			args: args{
				token: noPayloadToken(),
			},
			want:    nil,
			want1:   noPayloadToken(),
			wantErr: true,
		},
		{
			name: "Worked",
			JA:   NewJWTauth(&privKey.PublicKey),
			args: args{
				token: fineToken(),
			},
			want: &models.Payload{
				UID:        "mock_UID",
				Role:       0,
				Scope:      []string{"mock"},
				Refresh:    []byte("mock_refresh"),
				RefreshURL: "mock_refresh_url",
			},
			want1:   fineToken(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.JA.Validate(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"JWTauth.Validade() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if got == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf(
						"JWTauth.Validade() got = %v, want %v",
						got,
						tt.want,
					)
				}
			} else {
				if !reflect.DeepEqual(*got, *tt.want) {
					t.Errorf(
						"JWTauth.Validade() got = %v, want %v",
						got,
						tt.want,
					)
				}
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf(
					"JWTauth.Validade() got1 = %v, want %v",
					got1,
					tt.want1,
				)
			}
		})
	}
}

func TestJWTauth_Tokenize(t *testing.T) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	type args struct {
		load models.Payload
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
		handle  func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "Tokenize valid",
			args: args{
				load: models.Payload{
					UID:        "u000001",
					Scope:      []string{""},
					Role:       1,
					Refresh:    []byte("refreshtk"),
					RefreshURL: "http://refresh.token",
				},
			},
			want: []byte("mock_token"),
			handle: func(w http.ResponseWriter, r *http.Request) {
				token := models.JwtDO{
					Token: []byte("mock_token"),
				}
				rest.JSON(w, http.StatusOK, token)
			},
			wantErr: false,
		},
		{
			name: "Tokenize invalid UIDP response",
			args: args{
				load: models.Payload{
					UID:        "u000001",
					Scope:      []string{""},
					Role:       1,
					Refresh:    []byte("refreshtk"),
					RefreshURL: "http://refresh.token",
				},
			},
			want:    nil,
			wantErr: true,
			handle: func(w http.ResponseWriter, r *http.Request) {
				body := struct {
					Token bool `json:"token"`
				}{
					Token: true,
				}
				rest.JSON(w, http.StatusOK, body)
			},
		},
		{
			name: "Tokenize invalid",
			args: args{
				load: models.Payload{
					UID:        "u000001",
					Scope:      []string{""},
					Role:       1,
					Refresh:    []byte("refreshtk"),
					RefreshURL: "http://refresh.token",
				},
			},
			want: nil,
			handle: func(w http.ResponseWriter, r *http.Request) {
				err := ierrors.NewError().InternalServer().Message("error").Build()
				rest.ERROR(w, err)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(tt.handle)

			ts := httptest.NewServer(handler)
			os.Setenv("AUTH_PATH", ts.URL)
			defer ts.Close()

			JA := NewJWTauth(&privKey.PublicKey)
			got, err := JA.Tokenize(tt.args.load)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTauth.Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JWTauth.Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJWTauth_Refresh(t *testing.T) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	type args struct {
		token []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
		handle  func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "Tokenize valid",
			args: args{
				token: mockenize(models.Payload{
					UID:        "u000001",
					Scope:      []string{""},
					Role:       1,
					Refresh:    []byte("refreshtk"),
					RefreshURL: "http://refresh.token",
				}),
			},
			want: []byte("mock_token"),
			handle: func(w http.ResponseWriter, r *http.Request) {
				token := models.JwtDO{
					Token: []byte("mock_token"),
				}
				rest.JSON(w, http.StatusOK, token)
			},
			wantErr: false,
		},
		{
			name: "Tokenize invalid token payload",
			args: args{
				token: []byte("not_token"),
			},
			want:    nil,
			handle:  nil,
			wantErr: true,
		},
		{
			name: "Tokenize invalid UID response",
			args: args{
				token: mockenize(models.Payload{
					UID:        "u000001",
					Scope:      []string{""},
					Role:       1,
					Refresh:    []byte("refreshtk"),
					RefreshURL: "http://refresh.token",
				}),
			},
			want: nil,
			handle: func(w http.ResponseWriter, r *http.Request) {
				body := struct {
					Token bool `json:"token"`
				}{
					Token: true,
				}
				rest.JSON(w, http.StatusOK, body)
			},
			wantErr: true,
		},
		{
			name: "Tokenize invalid UID refresh",
			args: args{
				token: mockenize(models.Payload{
					UID:        "u000001",
					Scope:      []string{""},
					Role:       1,
					Refresh:    []byte("refreshtk"),
					RefreshURL: "http://refresh.token",
				}),
			},
			want: nil,
			handle: func(w http.ResponseWriter, r *http.Request) {
				err := ierrors.NewError().InternalServer().Message("error").Build()
				rest.ERROR(w, err)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(tt.handle)

			ts := httptest.NewServer(handler)
			os.Setenv("AUTH_PATH", ts.URL)
			defer ts.Close()

			JA := NewJWTauth(&privKey.PublicKey)
			got, err := JA.Refresh(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTauth.Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JWTauth.Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockenize(load models.Payload) []byte {
	token := jwt.New()
	token.Set(jwt.ExpirationKey, time.Now().Add(30*time.Minute))
	token.Set("payload", load)
	key, _ := rsa.GenerateKey(rand.Reader, 512)
	signed, _ := jwt.Sign(token, jwa.RS256, key)
	return signed
}
