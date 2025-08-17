package gcfg

import (
	"reflect"
	"testing"
)

type Config struct {
	Pnt    Point    `gcfg:"Point"`
	Set    bool     `gcfg:"set"`
	SecArr []SecArr `gcfg:"SecArr"`
}

type Point struct {
	X    int32    `gcfg:"x"`
	Y    int32    `gcfg:"y"`
	Z    int8     `gcfg:"z"`
	S    []int32  `gcfg:"s"`
	H    []string `gcfg:"h"`
	Name string   `gcfg:"name"`
}

type SecArr struct {
	Foo int32 `gcfg:"foo"`
}

func TestUnmarshal(t *testing.T) {
	input := `
Point {
	x = 3
	y = 1
	z = 4
	s = [1,2,3,4]
	h = ["h", "i"]
	name = "hello"
}

set = true

[SecArr] {
	foo = 3
}

[SecArr] {
	foo = 4
}
`
	var expectedCfg = Config{
		Pnt: Point{
			X:    3,
			Y:    1,
			Z:    4,
			S:    []int32{1, 2, 3, 4},
			H:    []string{"h", "i"},
			Name: "hello",
		},
		Set: true,
		SecArr: []SecArr{
			{
				Foo: 3,
			},
			{
				Foo: 4,
			},
		},
	}
	var cfg Config
	err := Unmarshal([]byte(input), &cfg)
	if err != nil || !reflect.DeepEqual(cfg, expectedCfg) {
		t.Errorf("Unmarshal=%v, %v want match for %v", cfg, err, expectedCfg)
	}
}
