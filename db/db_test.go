package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "SUCCESS",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init()
			if err != nil && !tt.wantErr {
				t.Error("Test Failed: {} Encountered error {}", tt.name, err)
			}
			assert.NotEmpty(t, got, "The returned context is not as expected")
		})
	}
}
