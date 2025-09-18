package fixedwidth

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type Processor[T any] struct {
	from      map[int]int
	to        map[int]int
	minLength int
}

func NewProcessor[T any]() (*Processor[T], error) {
	var v T
	FROM := make(map[int]int)
	TO := make(map[int]int)

	t := reflect.TypeOf(v)
	if t == nil {
		return nil, fmt.Errorf("%w: %T", ErrInvalidTarget, v)
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: %T", ErrInvalidTarget, v)
	}
	minLength := 0
	for i := range t.NumField() {
		tags := t.Field(i).Tag
		rawFrom := tags.Get("from")
		from, err := strconv.Atoi(rawFrom)
		if err != nil {
			return nil, fmt.Errorf("%w ('%s'): %w", ErrInvalidFrom, rawFrom, err)
		}
		rawTo := tags.Get("to")
		to, err := strconv.Atoi(rawTo)
		if err != nil {
			return nil, fmt.Errorf("%w ('%s'): %w", ErrInvalidTo, rawTo, err)
		}
		if from >= to {
			return nil, fmt.Errorf("%w: (%s:%s)", ErrInvalidFromTo, rawFrom, rawTo)
		}
		FROM[i] = from
		TO[i] = to
		if to > minLength {
			minLength = to
		}
	}
	return &Processor[T]{
		from:      FROM,
		to:        TO,
		minLength: minLength,
	}, nil
}

func (this *Processor[T]) Parse(line string, result *T) error {
	if len(line) < this.minLength {
		return fmt.Errorf("%w: line too short", ErrLineOverflow)
	}
	v := reflect.ValueOf(result).Elem()
	for i := range maps.Keys(this.from) {
		v.Field(i).Set(reflect.ValueOf(line[this.from[i]:this.to[i]]))
	}
	return nil
}
func (this *Processor[T]) Fprint(writer io.Writer, t T) (nn int, err error) {
	v := reflect.ValueOf(t)
	for _, i := range slices.Sorted(maps.Keys(this.from)) {
		length := this.to[i] - this.from[i]
		template := fmt.Sprintf("%%-%ds", length)
		n, err := fmt.Fprintf(writer, template, v.Field(i).String())
		nn += n
		if err != nil {
			return nn, err
		}
	}
	return nn, err
}
func (this *Processor[T]) Fprintln(writer io.Writer, v T) (nn int, err error) {
	nn, err = this.Fprint(writer, v)
	if err != nil {
		return nn, err
	}
	n, err := io.WriteString(writer, "\n")
	nn += n
	return nn, err
}
func (this *Processor[T]) Sprint(v T) string {
	var result strings.Builder
	_, _ = this.Fprint(&result, v)
	return result.String()
}
func (this *Processor[T]) Sprintln(v T) string {
	return this.Sprint(v) + "\n"
}

var (
	ErrInvalidTarget = errors.New("invalid target type")
	ErrInvalidFrom   = errors.New("invalid 'from' value")
	ErrInvalidTo     = errors.New("invalid 'to' value")
	ErrInvalidFromTo = errors.New("invalid from/to (from > to)")
	ErrLineOverflow  = errors.New("line overflow")
)
