package fixedwidth

import (
	"testing"

	"github.com/mdw-go/testing/should"
)

func TestOk(t *testing.T) {
	type Letters struct {
		A string `from:"0" to:"1"`
		B string `from:"1" to:"3"`
		C string `from:"3" to:"6"`
		D string `from:"6" to:"10"`
	}
	var letters Letters
	err := Unmarshal("a"+"bb"+"ccc"+"dddd", &letters)
	should.So(t, err, should.BeNil)
	should.So(t, letters, should.Equal, Letters{
		A: "a",
		B: "bb",
		C: "ccc",
		D: "dddd",
	})
}
func TestBlank(t *testing.T) {
	err := Unmarshal("", &struct{}{})
	should.So(t, err, should.WrapError, ErrEmptyLine)
}
func TestNilResultType(t *testing.T) {
	var v any
	err := Unmarshal("hi", v)
	should.So(t, err, should.WrapError, ErrInvalidTarget)
}
func TestInvalidResultType(t *testing.T) {
	err := Unmarshal("hi", make(chan int))
	should.So(t, err, should.WrapError, ErrInvalidTarget)
}
func TestLineOverflow(t *testing.T) {
	type Data struct {
		Overflows string `from:"0" to:"5"`
	}
	var data Data
	err := Unmarshal("1234", &data)
	should.So(t, err, should.WrapError, ErrLineOverflow)
}
func TestFromAboveTo(t *testing.T) {
	type Data struct {
		Field string `from:"1" to:"0"`
	}
	var data Data
	err := Unmarshal("1234", &data)
	should.So(t, err, should.WrapError, ErrInvalidFromTo)
}
func TestInvalidFrom(t *testing.T) {
	type Data struct {
		Field string `from:"NaN" to:"0"`
	}
	var data Data
	err := Unmarshal("1234", &data)
	should.So(t, err, should.WrapError, ErrInvalidFrom)
}
func TestInvalidTo(t *testing.T) {
	type Data struct {
		Field string `from:"0" to:"NaN"`
	}
	var data Data
	err := Unmarshal("1234", &data)
	should.So(t, err, should.WrapError, ErrInvalidTo)
}
