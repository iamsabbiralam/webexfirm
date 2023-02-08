package postgres

import (
	"context"
	"sort"
	"testing"

	"practice/webex/hrm/storage"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStorage_CreateCircularCategory(t *testing.T) {
	ts := newTestStorage(t)
	tests := []struct {
		name      string
		in        storage.CircularCategory
		want      *storage.CircularCategory
		multiWant []storage.CircularCategory
		wantErr   bool
	}{
		{
			name: "SUCCESS_CIRCULAR_CATEGORY",
			in: storage.CircularCategory{
				Name:        "test",
				Description: "test description",
				Status:      1,
				Position:    1,
				CRUDTimeDate: storage.CRUDTimeDate{
					CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
				},
			},
			want: &storage.CircularCategory{
				Name:        "test",
				Description: "test description",
				Status:      1,
				Position:    1,
				CRUDTimeDate: storage.CRUDTimeDate{
					CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
				},
			},
			multiWant: []storage.CircularCategory{
				{
					Name:        "test",
					Description: "test description",
					Status:      1,
					Position:    1,
					CRUDTimeDate: storage.CRUDTimeDate{
						CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
					},
				},
			},
		},
	}

	defer ts.db.Exec("DELETE FROM circular_categories")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ts.CreateCircularCategory(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateCircularCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.in.ID = got
			gotList, err := ts.ListCircularCategory(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListCircularCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			wantList := tt.multiWant
			sort.Slice(wantList, func(i, j int) bool {
				return wantList[i].ID < wantList[j].ID
			})
			sort.Slice(gotList, func(i, j int) bool {
				return gotList[i].ID < gotList[j].ID
			})
			for i, got := range gotList {
				tOps := []cmp.Option{
					cmpopts.IgnoreFields(storage.CircularCategory{}, "ID", "Count"),
					cmpopts.IgnoreFields(storage.CRUDTimeDate{}, "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy"),
				}

				if !cmp.Equal(got, wantList[i], tOps...) {
					t.Errorf("Diff: got -, want += %v", cmp.Diff(got, wantList[i]))
				}
			}

			result, err := ts.GetCircularCategory(context.Background(), tt.in.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetReturnCause() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tOps := []cmp.Option{
				cmpopts.IgnoreFields(storage.CircularCategory{}, "ID", "Count"),
				cmpopts.IgnoreFields(storage.CRUDTimeDate{}, "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy"),
			}
			if !cmp.Equal(result, tt.want, tOps...) {
				t.Errorf("Diff: got -, want += %v", cmp.Diff(result, tt.want))
			}

			resUp, err := ts.UpdateCircularCategory(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetCircularCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tOps = []cmp.Option{
				cmpopts.IgnoreFields(storage.CircularCategory{}, "ID", "Count"),
				cmpopts.IgnoreFields(storage.CRUDTimeDate{}, "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy"),
			}
			if !tt.wantErr && !cmp.Equal(tt.want, resUp, tOps...) {
				t.Errorf("Storage.UpdateCircularCategory() = + got, - want: %+v", cmp.Diff(tt.want, resUp))
			}

			tt.in = storage.CircularCategory{}
			tt.in.ID = got
			tt.in.DeletedBy = result.DeletedBy
			err = ts.DeleteCircularCategory(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteCircularCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
