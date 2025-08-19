package pair

import "testing"

var testPair = Pair[string, int]{
	First:  "hello",
	Second: 3,
}

func TestPair_String(t *testing.T) {
	expectedStr := "(hello, 3)"

	if testPair.String() != expectedStr {
		t.Errorf("Pair.String()=%s, wanted %s", testPair.String(), expectedStr)
	}
}

func TestPair_Values(t *testing.T) {
	str, integer := testPair.Values()

	expectedStr := "hello"
	expectedInteger := 3

	if str != expectedStr || integer != expectedInteger {
		t.Errorf("Pair.Values()=%s, %d, wanted %s, %d", str, integer, expectedStr, expectedInteger)
	}
}
