package parser

import (
	"fmt"
	"gcfg/lexer"
	"reflect"
	"testing"
)

func TestSimpleAssign(t *testing.T) {
	input := `
x = 3
y = 4.4
z = "hello"
b = true
c = false
d = nil
`
	expectedOutput := map[string]any{
		"x": 3,
		"y": 4.4,
		"z": "hello",
		"b": true,
		"c": false,
		"d": nil,
	}

	l := lexer.New([]byte(input))
	p := New(l)

	output, err := p.ParseFile()

	fmt.Println(output)

	if err != nil || !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("ParseFile=%v, %v, wanted match for %v", output, err, expectedOutput)
	}
}
