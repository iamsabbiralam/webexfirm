package postgres

import (
	"context"
	"testing"

	"practice/webex/hrm/storage"
)

func TestStorageGetUser(t *testing.T) {
	ts := newTestStorage(t)
	tests := []struct {
		name    string
		in      storage.User
		wantErr bool
	}{
		{
			name: "LOGIN_SUCCESS",
			in: storage.User{
				Email:    "test@email.com",
				Password: "test@1234",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ts.GetUser(context.TODO(), tt.in); (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
