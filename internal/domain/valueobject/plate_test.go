package valueobject_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewPlate_ShouldReturnPlateWhenOldFormatWithHyphen(t *testing.T) {
	plate, err := valueobject.NewPlate("ABC-1234")
	assert.NoError(t, err)
	assert.Equal(t, "ABC1234", plate.Value())
}

func TestNewPlate_ShouldReturnPlateWhenOldFormatWithoutHyphen(t *testing.T) {
	plate, err := valueobject.NewPlate("ABC1234")
	assert.NoError(t, err)
	assert.Equal(t, "ABC1234", plate.Value())
}

func TestNewPlate_ShouldReturnPlateWhenMercosulFormat(t *testing.T) {
	plate, err := valueobject.NewPlate("ABC1D23")
	assert.NoError(t, err)
	assert.Equal(t, "ABC1D23", plate.Value())
}

func TestNewPlate_ShouldReturnPlateWhenLowercase(t *testing.T) {
	plate, err := valueobject.NewPlate("abc-1234")
	assert.NoError(t, err)
	assert.Equal(t, "ABC1234", plate.Value())
}

func TestNewPlate_ShouldReturnErrorWhenInvalidFormat(t *testing.T) {
	_, err := valueobject.NewPlate("1234ABC")
	assert.ErrorIs(t, err, valueobject.ErrInvalidPlate)
}

func TestNewPlate_ShouldReturnErrorWhenTooShort(t *testing.T) {
	_, err := valueobject.NewPlate("ABC123")
	assert.ErrorIs(t, err, valueobject.ErrInvalidPlate)
}

func TestNewPlate_ShouldReturnErrorWhenEmpty(t *testing.T) {
	_, err := valueobject.NewPlate("")
	assert.ErrorIs(t, err, valueobject.ErrInvalidPlate)
}
