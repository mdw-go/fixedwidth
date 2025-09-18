package fixedwidth

import (
	"bytes"
	"testing"

	"github.com/mdw-go/testing/should"
)

type Letters struct {
	A string `from:"0" to:"1"`
	B string `from:"1" to:"3"`
	C string `from:"3" to:"6"`
	D string `from:"6" to:"10"`
}

func TestParse(t *testing.T) {
	processor, err := NewProcessor[Letters]()
	should.So(t, err, should.BeNil)
	var letters Letters
	input := "a" + "bb" + "ccc" + "dddd"
	err = processor.Parse(input, &letters)
	should.So(t, err, should.BeNil)
	should.So(t, letters, should.Equal, Letters{
		A: "a",
		B: "bb",
		C: "ccc",
		D: "dddd",
	})
	var buffer bytes.Buffer
	n, err := processor.Fprintln(&buffer, letters)
	should.So(t, err, should.BeNil)
	should.So(t, n, should.Equal, len(input+"\n"))
	should.So(t, buffer.String(), should.Equal, input+"\n")
	should.So(t, processor.Sprintln(letters), should.Equal, input+"\n")
}
