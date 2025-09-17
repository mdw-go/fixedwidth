package fixedwidth

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Marshal(v any) (string, error) {
	return "", nil // TODO: (requires putting all fields in order of from/to, and ensuring there aren't any overlapping fields)
}
func Unmarshal(line string, result any) error {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	if len(line) == 0 {
		return ErrEmptyLine
	}
	t := reflect.TypeOf(result)
	if t == nil {
		return fmt.Errorf("%w: %T", ErrInvalidTarget, result)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %T", ErrInvalidTarget, result)
	}
	v := reflect.ValueOf(result).Elem()
	for i := range t.NumField() {
		tags := t.Field(i).Tag
		rawFrom := tags.Get("from")
		from, err := strconv.Atoi(rawFrom)
		if err != nil {
			return fmt.Errorf("%w ('%s'): %w", ErrInvalidFrom, rawFrom, err)
		}
		rawTo := tags.Get("to")
		to, err := strconv.Atoi(rawTo)
		if err != nil {
			return fmt.Errorf("%w ('%s'): %w", ErrInvalidTo, rawTo, err)
		}
		if from >= to {
			return fmt.Errorf("%w: (%s:%s)", ErrInvalidFromTo, rawFrom, rawTo)
		}
		if to > len(line) {
			return fmt.Errorf("%w: (%s:%s)", ErrLineOverflow, rawFrom, rawTo)
		}
		v.Field(i).Set(reflect.ValueOf(line[from:to]))
	}
	return nil
}

var (
	ErrEmptyLine     = errors.New("empty line")
	ErrInvalidTarget = errors.New("invalid target type")
	ErrInvalidFrom   = errors.New("invalid 'from' value")
	ErrInvalidTo     = errors.New("invalid 'to' value")
	ErrInvalidFromTo = errors.New("invalid from/to (from > to)")
	ErrLineOverflow  = errors.New("line overflow")
)
