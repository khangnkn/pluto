package taskapi

import (
	"testing"

	"github.com/sebdah/goldie/v2"
)

func Test_pushTaskMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pushTaskMessage()
			gd := goldie.New(t)
			gd.AssertJson(t, t.Name(), got)
		})
	}
}
