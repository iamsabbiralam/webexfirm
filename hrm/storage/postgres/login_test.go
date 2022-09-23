package postgres

import (
	"context"
	"personal/webex/hrm/storage"
	"testing"
	"time"
)

func TestStorageLogin(t *testing.T) {
	ts := newTestStorage(t)
	tests := []struct {
		name    string
		in      string
		want    storage.SignUP
		wantErr bool
	}{
		{
			name: "LOGIN_SUCCESS",
			in:   "test@email.com",
			want: storage.SignUP{
				ID:        "user-ID",
				FirstName: "Sabbir",
				LastName:  "Alam",
				Username:  "iamsabbiralam",
				Email:     "sabbir@gmail.com",
				Image:     "image.jpg",
				Phone:     "01715039303",
				Password:  "sabbir",
				Gender:    1,
				DOB:       time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
				Status:    1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ts.Login(context.TODO(), tt.in); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
