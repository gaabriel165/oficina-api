package valueobject

import (
	"errors"
	"strconv"

	"github.com/gabrielcamargo/oficina-api/internal/utils"
)

var ErrInvalidCPF = errors.New("invalid CPF")

type CPF struct {
	value string
}

func NewCPF(value string) (CPF, error) {
	cleaned := utils.CleanNumericString(value)
	if !isValidCPF(cleaned) {
		return CPF{}, ErrInvalidCPF
	}
	return CPF{value: cleaned}, nil
}

func (c CPF) Value() string {
	return c.value
}

func (c CPF) String() string {
	return c.value
}

func isValidCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}

	if allDigitsSame(cpf) {
		return false
	}

	if !validateCPFDigit(cpf, 9) {
		return false
	}

	return validateCPFDigit(cpf, 10)
}

func allDigitsSame(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func validateCPFDigit(cpf string, position int) bool {
	sum := 0
	for i := 0; i < position; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (position + 1 - i)
	}

	remainder := sum % 11
	expected := 0
	if remainder >= 2 {
		expected = 11 - remainder
	}

	actual, _ := strconv.Atoi(string(cpf[position]))
	return actual == expected
}
