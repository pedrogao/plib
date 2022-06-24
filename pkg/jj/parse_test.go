package jj

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pedrogao/plib/pkg/common"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    Elem
		wantErr bool
	}{
		{
			name:    "[1, 2, true]",
			source:  "[1, 2, true]",
			want:    NewArrayElem([]Elem{NewLiteralElem(1.0), NewLiteralElem(2.0), NewLiteralElem(true)}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.source)
			tokens, err := scanner.ScanTokens()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			p := NewParser(tokens)
			got, err := p.Parse()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Parse1(t *testing.T) {
	assert := assert.New(t)
	s := testdata("./testdata/1.json")
	scanner := NewScanner(s)
	tokens, err := scanner.ScanTokens()
	assert.Nil(err)
	p := NewParser(tokens)
	got, err := p.Parse()
	assert.Nil(err)
	t.Logf("%+v\n", got)
	assert.Equal(got.Type(), ObjectType)
}

func TestParser_Parse2(t *testing.T) {
	assert := assert.New(t)
	s := testdata("./testdata/2.json")
	scanner := NewScanner(s)
	tokens, err := scanner.ScanTokens()
	assert.Nil(err)
	p := NewParser(tokens)
	got, err := p.Parse()
	assert.Nil(err)
	t.Logf("%+v\n", got)
	assert.Equal(got.Type(), ObjectType)
}

func TestParser_Parse3(t *testing.T) {
	assert := assert.New(t)
	s := testdata("./testdata/3.json")
	scanner := NewScanner(s)
	tokens, err := scanner.ScanTokens()
	assert.Nil(err)
	p := NewParser(tokens)
	got, err := p.Parse()
	assert.Nil(err)
	t.Logf("%+v\n", got)
	assert.Equal(got.Type(), ArrayType)
	assert.Equal(got.AsArray()[0].AsObject()["id"].AsString(), "Open")
	assert.Equal(got.Get(0, "id").AsString(), "Open")
}

func TestParser_Parse4(t *testing.T) {
	assert := assert.New(t)
	s := testdata("./testdata/4.json")
	scanner := NewScanner(s)
	tokens, err := scanner.ScanTokens()
	assert.Nil(err)
	p := NewParser(tokens)
	got, err := p.Parse()
	assert.Nil(err)
	t.Logf("%+v\n", got)
	assert.Equal(got.Type(), ArrayType)
	assert.Equal(got.AsArray()[0].AsNumber(), 1.0)
	assert.Equal(got.Get(1).AsNumber(), 2.0)
}

// "./testdata/1.json"
func testdata(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return common.Bytes2String(data)
}
