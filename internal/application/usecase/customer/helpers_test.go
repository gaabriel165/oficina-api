package customer_test

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

func makeCustomer() (*entity.Customer, error) {
	return entity.NewCustomer("John Doe", "529.982.247-25", "11999999999", "john@email.com")
}
