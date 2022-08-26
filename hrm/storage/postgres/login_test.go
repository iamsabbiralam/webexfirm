package postgres

import (
	"context"
	"personal/webex/hrm/storage"
	"testing"
)

func TestStorageLogin(t *testing.T) {
	ts := newTestStorage(t)
	tests := []struct {
		name    string
		in      storage.SignUP
		wantErr bool
	}{
		{
			name: "LOGIN_SUCCESS",
			in: storage.SignUP{
				Email:    "test@email.com",
				Password: "test@1234",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := ts.Login(context.TODO(), tt.in); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
