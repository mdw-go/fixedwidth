package fixedwidth

import (
	"testing"

	"github.com/mdw-go/testing/should"
)

func TestDecoder(t *testing.T) {
	type Letters struct {
		A string `from:"0" to:"1"`
		B string `from:"1" to:"3"`
		C string `from:"3" to:"6"`
		D string `from:"6" to:"10"`
	}
	decoder, err := NewDecoder[Letters]()
	should.So(t, err, should.BeNil)
	var letters Letters
	err = decoder.Decode("a"+"bb"+"ccc"+"dddd", &letters)
	should.So(t, err, should.BeNil)
	should.So(t, letters, should.Equal, Letters{
		A: "a",
		B: "bb",
		C: "ccc",
		D: "dddd",
	})
}
