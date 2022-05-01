package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ctx := Context{}
	tests := []struct {
		name    string
		want    *Context
		wantErr bool
	}{
		{
			name:    "SUCCESS",
			want:    &Context{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctx.Init()
			if err != nil && !tt.wantErr {
				t.Error("Test Failed: {} Encountered error {}", tt.name, err)
			}
			assert.NotEmpty(t, got, "The returned context is not as expected")
		})
	}
}
