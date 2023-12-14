package pkg

import (
	"Practica/internal/domain"
	"encoding/json"
	"io/ioutil"
)

func FullfilDB(filePath string) []domain.Product {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var products []domain.Product
	err = json.Unmarshal(file, &products)
	if err != nil {
		panic(err)
	}

	return products
}
