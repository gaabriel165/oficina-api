package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidPlate = errors.New("invalid vehicle plate")

var (
	oldFormatPlate      = regexp.MustCompile(`^[A-Z]{3}\d{4}$`)
	mercosulFormatPlate = regexp.MustCompile(`^[A-Z]{3}\d[A-Z]\d{2}$`)
)

type Plate struct {
	value string
}

func NewPlate(value string) (Plate, error) {
	normalized := strings.ToUpper(strings.ReplaceAll(value, "-", ""))
	if !oldFormatPlate.MatchString(normalized) && !mercosulFormatPlate.MatchString(normalized) {
		return Plate{}, ErrInvalidPlate
	}
	return Plate{value: normalized}, nil
}

func (p Plate) Value() string {
	return p.value
}

func (p Plate) String() string {
	return p.value
}
