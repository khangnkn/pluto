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
				Title: "0",
			}, {
				Title: "1",
			}, {
				Title: "2",
			}, {
				Title: "3",
			}, {
				Title: "4",
			}, {
				Title: "5",
			}, {
				Title: "6",
			}}
	)
	type args struct {
		imgs   []image.Image
		cursor int
		s      int
	}
	tests := []struct {
		name string
		args args
		want []image.Image
		cur  int
	}{
		{
			name: "tc1",
			args: args{
				imgs:   testData,
				cursor: 0,
				s:      3,
			},
			want: testData[:3],
			cur:  3,
		},
		{
			name: "tc2",
			args: args{
				imgs:   testData,
				cursor: 2,
				s:      3,
			},
			want: testData[2:5],
			cur:  5,
		},
		{
			name: "tc3",
			args: args{
				imgs:   testData,
				cursor: 4,
				s:      3,
			},
			want: testData[4:7],
			cur:  7,
		},
		{
			name: "tc4",
			args: args{
				imgs:   testData,
				cursor: 4,
				s:      6,
			},
			want: append(testData[4:7], testData[:3]...),
			cur:  3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncate(tt.args.imgs, &tt.args.cursor, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("truncate() = %+v, want %+v", got, tt.want)
			} else {
				t.Logf("got %+v", got)
			}
		})
		if tt.args.cursor != tt.cur {
			t.Errorf("got %d, want %d", tt.args.cursor, tt.cur)
		}
	}
}
