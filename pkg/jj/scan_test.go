package jj

import (
	"reflect"
	"testing"
)

func TestScanner_ScanTokens(t *testing.T) {
	type fields struct {
		source string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*Token
		wantErr bool
	}{
		{
			name: "simple scan",
			fields: fields{
				source: "[1, true, \"pedro\"]",
			},
			want: []*Token{
				{
					Type:  LeftSquare,
					Raw:   "[",
					Value: nil,
					Line:  1,
				},
				{
					Type:  Number,
					Raw:   "1",
					Value: 1.0,
					Line:  1,
				},
				{
					Type:  Comma,
					Raw:   ",",
					Value: nil,
					Line:  1,
				},
				{
					Type:  True,
					Raw:   "true",
					Value: true,
					Line:  1,
				},
				{
					Type:  Comma,
					Raw:   ",",
					Value: nil,
					Line:  1,
				},
				{
					Type:  String,
					Raw:   "\"pedro\"",
					Value: "pedro",
					Line:  1,
				},
				{
					Type:  RightSquare,
					Raw:   "]",
					Value: nil,
					Line:  1,
				},
				{
					Type:  Eof,
					Raw:   "",
					Value: nil,
					Line:  1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.fields.source)
			got, err := s.ScanTokens()
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanTokens() got = %v, want %v", got, tt.want)
			}
		})
	}
}
