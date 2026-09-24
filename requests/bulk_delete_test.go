package requests

import (
	"testing"
)

func TestBulkDeleteRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     BulkDeleteRequest
		wantErr bool
	}{
		{
			name:    "valid IDs",
			req:     BulkDeleteRequest{IDs: []uint{1, 2, 5, 10}},
			wantErr: false,
		},
		{
			name:    "single valid ID",
			req:     BulkDeleteRequest{IDs: []uint{1}},
			wantErr: false,
		},
		{
			name:    "empty IDs",
			req:     BulkDeleteRequest{IDs: []uint{}},
			wantErr: true,
		},
		{
			name:    "contains zero ID",
			req:     BulkDeleteRequest{IDs: []uint{1, 0, 5}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
