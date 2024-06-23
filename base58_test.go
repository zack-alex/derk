package derk

import (
	"testing"
)

func TestB58Encode(t *testing.T) {
	tests := []struct {
		given  []byte
		expect string
	}{
		{[]byte("Hello World!"), "2NEpo7TZRRrLZSi2U"},
		{[]byte(""), ""},
		{[]byte{0x00}, "1"},
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00}, "11111"},
		{[]byte{0x00, 0x00, 0x01}, "112"},
		{[]byte{0x01, 0x00, 0x00}, "LUw"},
	}

	for _, test := range tests {
		result := base58Encode(test.given)
		if result != test.expect {
			t.Errorf("b58encode(%q) = %s; expect %s", test.given, result, test.expect)
		}
	}
}
