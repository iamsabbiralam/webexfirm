package postgres

import (
	"context"
	"sort"
	"testing"

	"practice/webex/hrm/storage"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStorage_CreateJobType(t *testing.T) {
	ts := newTestStorage(t)
	tests := []struct {
		name      string
		in        storage.JobTypes
		want      *storage.JobTypes
		multiWant []storage.JobTypes
		wantErr   bool
	}{
		{
			name: "SUCCESS_JOB_TYPES",
			in: storage.JobTypes{
				Name:     "test",
				Status:   1,
				Position: 1,
				CRUDTimeDate: storage.CRUDTimeDate{
					CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
				},
			},
			want: &storage.JobTypes{
				Name:     "test",
				Status:   1,
				Position: 1,
				CRUDTimeDate: storage.CRUDTimeDate{
					CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
				},
			},
			multiWant: []storage.JobTypes{
				{
					Name:     "test",
					Status:   1,
					Position: 1,
					CRUDTimeDate: storage.CRUDTimeDate{
						CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ts.CreateJobType(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateJobType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.in.ID = got
			gotList, err := ts.ListJobTypes(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListJobTypes() error = %v, wantErr %v", err, tt.wantErr)
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
					cmpopts.IgnoreFields(storage.JobTypes{}, "ID", "Count"),
					cmpopts.IgnoreFields(storage.CRUDTimeDate{}, "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy"),
				}

				if !cmp.Equal(got, wantList[i], tOps...) {
					t.Errorf("Diff: got -, want += %v", cmp.Diff(got, wantList[i]))
				}
			}

			result, err := ts.GetJobType(context.Background(), tt.in.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetJobType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tOps := []cmp.Option{
				cmpopts.IgnoreFields(storage.JobTypes{}, "ID", "Count"),
				cmpopts.IgnoreFields(storage.CRUDTimeDate{}, "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy"),
			}
			if !cmp.Equal(result, tt.want, tOps...) {
				t.Errorf("Diff: got -, want += %v", cmp.Diff(result, tt.want))
			}

			resUp, err := ts.UpdateJobType(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.UpdateJobType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tOps = []cmp.Option{
				cmpopts.IgnoreFields(storage.JobTypes{}, "ID", "Count"),
				cmpopts.IgnoreFields(storage.CRUDTimeDate{}, "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy"),
			}
			if !tt.wantErr && !cmp.Equal(tt.want, resUp, tOps...) {
				t.Errorf("Storage.UpdateJobType() = + got, - want: %+v", cmp.Diff(tt.want, resUp))
			}

			tt.in = storage.JobTypes{}
			tt.in.ID = got
			tt.in.DeletedBy = result.DeletedBy
			err = ts.DeleteJobType(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteJobType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
