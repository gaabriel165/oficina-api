package valueobject

import (
	"errors"
	"strconv"

	"github.com/gabrielcamargo/oficina-api/internal/utils"
)

var ErrInvalidCNPJ = errors.New("invalid CNPJ")

type CNPJ struct {
	value string
}

func NewCNPJ(value string) (CNPJ, error) {
	cleaned := utils.CleanNumericString(value)
	if !isValidCNPJ(cleaned) {
		return CNPJ{}, ErrInvalidCNPJ
	}
	return CNPJ{value: cleaned}, nil
}

func (c CNPJ) Value() string {
	return c.value
}

func (c CNPJ) String() string {
	return c.value
}

func isValidCNPJ(cnpj string) bool {
	if len(cnpj) != 14 {
		return false
	}

	if allDigitsSame(cnpj) {
		return false
	}

	if !validateCNPJDigit(cnpj, []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}, 12) {
		return false
	}

	return validateCNPJDigit(cnpj, []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}, 13)
}

func validateCNPJDigit(cnpj string, weights []int, position int) bool {
	sum := 0
	for i := 0; i < position; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights[i]
	}

	remainder := sum % 11
	expected := 0
	if remainder >= 2 {
		expected = 11 - remainder
	}

	actual, _ := strconv.Atoi(string(cnpj[position]))
	return actual == expected
}
