package postgres

import (
	"context"
	"personal/webex/hrm/storage"
	"testing"
	"time"
)

func TestStorage_SignUP(t *testing.T) {
	ts := newTestStorage(t)
	tests := []struct {
		name    string
		in      storage.SignUP
		want    string
		wantErr bool
	}{
		{
			name: "SIGNUP_SUCCESS",
			in: storage.SignUP{
				FirstName: "Sabbir",
				LastName:  "Alam",
				Username:  "iamsabbiralam",
				Email:     "sabbir@webex.com",
				Image:     "image.jpg",
				Phone:     "01715039303",
				Password:  "sabbir007",
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
			_, err := ts.SignUP(context.TODO(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.SignUP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
