package pair

import "fmt"

type Pair[T, U any] struct {
	First  T
	Second U
}

func (p Pair[T, U]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

func (p Pair[T, U]) Values() (T, U) {
	return p.First, p.Second
}
