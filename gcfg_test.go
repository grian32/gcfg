package gcfg

import (
	"reflect"
	"testing"

	"github.com/grian32/gcfg/pair"
)

type Config struct {
	Pnt    Point    `gcfg:"Point"`
	Set    bool     `gcfg:"set"`
	SecArr []SecArr `gcfg:"SecArr"`
}

type Point struct {
	X          int32                    `gcfg:"x"`
	Y          int32                    `gcfg:"y"`
	Z          uint8                    `gcfg:"z"`
	S          []int32                  `gcfg:"s"`
	L          []uint32                 `gcfg:"l"`
	Empty      []string                 `gcfg:"empty"`
	Ab         pair.Pair[int32, string] `gcfg:"ab"`
	H          []string                 `gcfg:"h"`
	Name       string                   `gcfg:"name"`
	Underscore string                   `gcfg:"und_erscore"`
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
	l = [1,2,3,4]
	empty = []
	ab = (1, "hi")
	h = ["h", "i"]
	name = "hello"
	und_erscore = "hi"
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
			X:          3,
			Y:          1,
			Z:          4,
			S:          []int32{1, 2, 3, 4},
			L:          []uint32{1, 2, 3, 4},
			Empty:      []string{},
			Ab:         pair.Pair[int32, string]{First: 1, Second: "hi"},
			H:          []string{"h", "i"},
			Name:       "hello",
			Underscore: "hi",
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
