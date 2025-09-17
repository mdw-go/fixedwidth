package fixedwidth

import (
	"fmt"
	"maps"
	"reflect"
	"strconv"
)

type Decoder[T any] struct {
	from      map[int]int
	to        map[int]int
	minLength int
}

func NewDecoder[T any]() (*Decoder[T], error) {
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

	return &Decoder[T]{
		from:      FROM,
		to:        TO,
		minLength: minLength,
	}, nil
}

func (this *Decoder[T]) Decode(line string, result *T) error {
	if len(line) < this.minLength {
		return fmt.Errorf("%w: line too short", ErrLineOverflow)
	}
	v := reflect.ValueOf(result).Elem()
	for i := range maps.Keys(this.from) {
		v.Field(i).Set(reflect.ValueOf(line[this.from[i]:this.to[i]]))
	}
	return nil
}
