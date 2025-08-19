package parser

import (
	"reflect"
	"testing"

	"github.com/grian32/gcfg/lexer"
	"github.com/grian32/gcfg/pair"
)

func TestParseFile(t *testing.T) {
	input := `
x = 3
y = 4.4
z = "hello"
b = true
c = false
d = nil
h = (2, 2)
m = [1,2,3,4,5]

Sec {
	b = 4
	hi = true
}

[SecArr] {
	foo = 4
}

[SecArr] {
	foo = 5
}
`
	// ints are outputted as strings by the parser and converted at reflection for easier bounds checking
	expectedOutput := map[string]any{
		"x": "3",
		"y": 4.4,
		"z": "hello",
		"b": true,
		"c": false,
		"d": nil,
		"h": pair.Pair[any, any]{
			First:  "2",
			Second: "2",
		},
		"m": []any{"1", "2", "3", "4", "5"},
		"Sec": map[string]any{
			"b":  "4",
			"hi": true,
		},
		"SecArr": []map[string]any{
			{
				"foo": "4",
			},
			{
				"foo": "5",
			},
		},
	}

	l := lexer.New([]byte(input))
	p := New(l)

	output, err := p.ParseFile()

	if err != nil || !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("ParseFile=%v, %v, wanted match for %v", output, err, expectedOutput)
	}
}
