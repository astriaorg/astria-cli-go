package keys

import (
	"crypto/ed25519"
	"encoding/hex"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDecryptKeyfile(t *testing.T) {
	// NOTE - this is a test private key! don't use for anything real
	privkeyBytes, _ := hex.DecodeString("158fb2953ecb5a4fd416ec345df586d88ed7494e09075e5cf872337eede03424")
	privkey := ed25519.NewKeyFromSeed(privkeyBytes)

	type args struct {
		keyfile  string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    ed25519.PrivateKey
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				keyfile:  "testdata/keystore/empty",
				password: "foobar",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "keyfile doesn't exist",
			args: args{
				keyfile:  "testdata/keystore/nofile",
				password: "foobar",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				keyfile:  "testdata/keystore/UTC--2024-04-25T13:47:31-06:00--b1f11f673cfd6aa4cc69b25f7f59bc89bccc62f3",
				password: "foobar",
			},
			want:    privkey,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptKeyfile(tt.args.keyfile, tt.args.password)
			if tt.wantErr && err == nil {
				t.Errorf("DecryptPrivateKey() error %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptPrivateKey() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveKeyfilePath(t *testing.T) {
	tests := []struct {
		name    string
		keydir  string
		want    string
		wantErr bool
	}{
		{
			name:    "directory",
			keydir:  "testdata/keystore",
			want:    "UTC--2024-04-25T13:47:31-06:00--b1f11f673cfd6aa4cc69b25f7f59bc89bccc62f3",
			wantErr: false,
		},
		{
			name:    "notfound",
			keydir:  "testdata/keystore/null",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveKeyfilePath(tt.keydir)
			if tt.name != "notfound" {
				got = filepath.Base(got)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveKeyfilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveKeyfilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
