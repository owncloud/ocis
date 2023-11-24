package preprocessor

import (
	"bytes"
	"image"
	"io"
	"reflect"

	"testing"
)

// implement tests for the following structs and functions:
// ImageDecoder
// GifDecoder
// GgsDecoder

// TestImageDecoderConvert tests the function Convert of the struct ImageDecoder
func TestImageDecoderConvert(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		i       ImageDecoder
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test successful convert image",
			args: args{
				r: bytes.NewReader([]byte("")),
			},
			want:    image.Image{},
			wantErr: false,
		},
		{
			name: "Test unsuccessful convert image",
			args: args{
				r: bytes.NewReader([]byte("")),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		i := ImageDecoder{}
		got, err := i.Convert(tt.args.r)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ImageDecoder.Convert() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
			t.Errorf("%q. ImageDecoder.Convert() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
