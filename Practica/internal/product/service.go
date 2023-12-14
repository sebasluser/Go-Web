package product

import (
	"Practica/internal/domain"
	"fmt"
	"time"
)

type ProductService struct {
	repository ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return ProductService{repo}
}

func (s *ProductService) GetAllProducts() []domain.Product {
	return s.repository.GetAllProducts()
}

func (s *ProductService) GetById(id int) domain.Product {
	return s.repository.GetById(id)
}

func (s *ProductService) GetByPrice(price float64) []domain.Product {
	products := s.repository.GetAllProducts()

	result := make([]domain.Product, 0)
	for _, product := range products {
		if product.Price > price {
			result = append(result, product)
		}
	}
	return result
}

func (s *ProductService) ValidateProductFields(product *domain.Product) error {
	if product.Name == "" || product.Quantity == 0 || product.CodeValue == "" ||
		product.Expiration == "" || product.Price == 0.0 {
		return fmt.Errorf("all fields except is_published must be provided")
	}
	return nil
}

func (s *ProductService) ValidateDateFormat(date string) error {
	_, err := time.Parse("02/01/2006", date)
	if err != nil {
		return fmt.Errorf("invalid date format for expiration")
	}
	return nil
}

func (s *ProductService) ValidateUniqueCodeValue(products []domain.Product, codeValue string) error {
	for _, existingProduct := range products {
		if existingProduct.CodeValue == codeValue {
			return fmt.Errorf("code_value must be unique")
		}
	}
	return nil
}

func (s *ProductService) AddProduct(newProduct domain.Product) {
	s.repository.AddProduct(newProduct)
}

func (s *ProductService) UpdateProduct(updatedProduct domain.Product) error {

	if err := s.ValidateProductFields(&updatedProduct); err != nil {
		return err
	}

	if err := s.ValidateDateFormat(updatedProduct.Expiration); err != nil {
		return err
	}

	if err := s.ValidateUniqueCodeValue(s.repository.GetAllProducts(), updatedProduct.CodeValue); err != nil {
		return err
	}

	return s.repository.UpdateProduct(updatedProduct)
}

func (s *ProductService) DeleteProduct(id int) error {
	if id < 1 || id > len(s.repository.GetAllProducts()) {
		return fmt.Errorf("product not found")
	}
	s.repository.DeleteProduct(id)
	return nil
}
