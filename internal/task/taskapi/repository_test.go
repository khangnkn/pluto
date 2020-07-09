package taskapi

import (
	"reflect"
	"testing"

	"github.com/nkhang/pluto/internal/image"

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

func Test_truncate(t *testing.T) {
	var (
		testData = []image.Image{
			{
				Title: "1",
			}, {
				Title: "2",
			}, {
				Title: "3",
			}, {
				Title: "4",
			}}
	)
	type args struct {
		imgs   []image.Image
		cursor *int
		s      int
	}
	tests := []struct {
		name string
		args args
		want []image.Image
	}{
		{
			name: "tc1",
			args: args{
				imgs:   testData,
				cursor: nil,
				s:      0,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncate(tt.args.imgs, tt.args.cursor, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}
