package encoder

import (
	"URLShortener/internal/core/encoder"
	"fmt"
	"testing"
)
import "math/rand"

func TestBase63EncodeSimple(t *testing.T) {
	test1 := int64(5)
	encoded := encoder.Base63Encode(test1)
	if encoded != "F" {
		t.Errorf("TestBase63Encode(%d) = %s, want %s", test1, encoded, "F")
	}

	test2 := int64(0)
	expected := "A"
	encoded2 := encoder.Base63Encode(test2)
	if encoded2 != expected {
		t.Errorf("TestBase63Encode(%d) = %s, want %s", test2, fmt.Sprintf("%q", encoded2), fmt.Sprintf("%q", expected))
	}

	test3 := int64(62)
	encoded3 := encoder.Base63Encode(test3)
	if encoded3 != "+" {
		t.Errorf("TestBase63Encode(%d) = %s, want %s", test3, encoded3, "+")
	}

	test4 := "4"
	decoded := encoder.Base63Decode(test4)
	if decoded != int64(56) {
		t.Errorf("TestBase63Decode(%s) = %d, want %d", test4, decoded, 56)
	}

	test5 := "aQ"
	decoded2 := encoder.Base63Decode(test5)
	if decoded2 != int64(1654) {
		t.Errorf("TestBase63Decode(%s) = %d, want %d", test5, decoded2, 1654)
	}
}

func TestBase63EncodeLong(t *testing.T) {
	for i := 0; i < 1000; i++ {
		testInt := rand.Int63()
		encoded := encoder.Base63Encode(testInt)
		decoded := encoder.Base63Decode(encoded)

		if testInt != decoded {
			t.Errorf("TestBase63Encode(%d) = %d, want %d", testInt, decoded, testInt)
			break
		}
	}
}
