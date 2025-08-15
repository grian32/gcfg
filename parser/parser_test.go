package parser

import (
	"fmt"
	"gcfg/lexer"
	"reflect"
	"testing"
)

func TestParseFile(t *testing.T) {
	input := `
x = 3
y = 4.4
z = "hello"
b = true
c = false
d = nil

Sec {
	b = 4
	hi = true
}
`
	expectedOutput := map[string]any{
		"x": 3,
		"y": 4.4,
		"z": "hello",
		"b": true,
		"c": false,
		"d": nil,
		"Sec": map[string]any{
			"b":  4,
			"hi": true,
		},
	}

	l := lexer.New([]byte(input))
	p := New(l)

	output, err := p.ParseFile()

	fmt.Println(output)

	if err != nil || !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("ParseFile=%v, %v, wanted match for %v", output, err, expectedOutput)
	}
}
