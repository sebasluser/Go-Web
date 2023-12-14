package product

import (
	"Practica/internal/domain"
	"fmt"
)

type ProductRepository struct {
	productsDB []domain.Product
}

func NewProductRepository(prods []domain.Product) ProductRepository {
	return ProductRepository{prods}
}

func (r ProductRepository) GetAllProducts() []domain.Product {
	return r.productsDB
}

func (r ProductRepository) GetById(id int) domain.Product {
	return r.productsDB[id-1]
}

func (r *ProductRepository) AddProduct(newProduct domain.Product) {
	r.productsDB = append(r.productsDB, newProduct)
}

func (r *ProductRepository) UpdateProduct(updatedProduct domain.Product) error {
	for i, existingProduct := range r.productsDB {
		if existingProduct.ID == updatedProduct.ID {
			r.productsDB[i] = updatedProduct
			return nil
		}
	}
	return fmt.Errorf("product not found")
}

// product/repository.go

func (r *ProductRepository) DeleteProduct(id int) {
	updatedProducts := make([]domain.Product, 0)

	for _, existingProduct := range r.productsDB {
		if existingProduct.ID != id {
			updatedProducts = append(updatedProducts, existingProduct)
		}
	}

	r.productsDB = updatedProducts
}
