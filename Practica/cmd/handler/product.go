package handler

import (
	"Practica/internal/domain"
	"Practica/internal/product"
	"Practica/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductRouter struct {
	productGroup *gin.RouterGroup
	service      product.ProductService
	repository   product.ProductRepository
}

func NewProductRouter(g *gin.RouterGroup) ProductRouter {
	slice := pkg.FullfilDB("../../products.json")
	repo := product.NewProductRepository(slice)

	serv := product.NewProductService(repo)

	return ProductRouter{productGroup: g, service: serv}

}

func (r *ProductRouter) ProductRoutes() {
	r.productGroup.GET("/ping", r.Ping())
	r.productGroup.GET("/getAll", r.GetAllProducts())
	r.productGroup.GET("/:id", r.GetById())
	r.productGroup.GET("/search", r.GetByPrice())
	r.productGroup.POST("/", r.AddProduct())
	r.productGroup.PUT("/", r.UpdateProduct())
	r.productGroup.PATCH("/", r.UpdateProduct())
	r.productGroup.DELETE("/:id", r.DeleteProduct())
}

func (r *ProductRouter) Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Pong")
	}
}

func (r *ProductRouter) GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := r.service.GetAllProducts()
		ctx.JSON(http.StatusOK, data)
	}
}

func (r *ProductRouter) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
			return
		}
		data := r.service.GetById(id)
		ctx.JSON(http.StatusOK, data)
	}
}

func (r *ProductRouter) GetByPrice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priceGt parameter"})
			return
		}
		data := r.service.GetByPrice(price)
		ctx.JSON(http.StatusOK, data)
	}
}

func (r *ProductRouter) AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProduct domain.Product

		if err := c.BindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if err := r.service.ValidateProductFields(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := r.service.ValidateDateFormat(newProduct.Expiration); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := r.service.ValidateUniqueCodeValue(r.repository.GetAllProducts(), newProduct.CodeValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newProduct.ID = len(r.repository.GetAllProducts()) + 1
		r.service.AddProduct(newProduct)

		c.IndentedJSON(http.StatusCreated, newProduct)
	}
}

func (r *ProductRouter) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updatedProduct domain.Product

		if err := c.BindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Validate and update the product
		if err := r.service.UpdateProduct(updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, updatedProduct)
	}
}

func (r *ProductRouter) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
			return
		}

		if err := r.service.DeleteProduct(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}
