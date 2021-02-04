package core

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type productService struct {
	repo     ProductRepo
	auth     JWTService
	cache	 Cache
	validate *validator.Validate
	publisherQ Publisher
}

func (p *productService) Authorize(s string) (map[string]interface{}, error) {
	return p.auth.ValidateToken(s)
}

func (p *productService) AllProducts(ctx context.Context) (Products, error) {
	var products Products
	err := p.cache.GetJSON(ctx, "all_products", &products)
	if err != nil {
		products, err = p.repo.All(ctx)
		if err != nil {
			return nil, err
		}
		p.cache.SetJSON(ctx, "all_products", &products, 0)
	}
	return products, nil
}

func (p *productService) ProductById(ctx context.Context, i int) (Product, error) {
	var product Product
	err := p.cache.GetJSON(ctx, "product_"+strconv.Itoa(i), &product)
	if err != nil {
		return p.repo.ByID(ctx, i)
	}
	return product, nil
}

func (p *productService) AddProduct(ctx context.Context, product *Product) (id int, err error) {

	if err = p.validate.Struct(product); err != nil {
		return
	}
	if id, err = p.repo.Add(ctx, product); err != nil {
		return
	}
	p.cache.SetJSON(ctx,  "product_"+strconv.Itoa(id), &product, 0)
	return
}

func (p *productService) EditProduct(ctx context.Context, product *Product) error {
	err := p.repo.Edit(ctx, product)
	if err != nil {
		return err
	}
	if product.Price > 0 {

		msg := []byte(fmt.Sprintf("%d|%f", product.Id, product.Price))
		p.publisherQ.Publish(ctx, &Message{
			Topic: "product",
			Msg:  msg,
		})
	}
	return p.cache.SetJSON(ctx,  "product_"+strconv.Itoa(product.Id), &product, 0)
}

func NewService(pr ProductRepo, cache Cache, auth JWTService, publisher Publisher) ProductService {

	val := validator.New()
	val.RegisterValidation("sku", validateSKU)

	return &productService{
		repo: pr,
		cache: cache,
		auth: auth,
		validate: val,
		publisherQ: publisher,
	}
}

func validateSKU(fl validator.FieldLevel) bool {
	// SKU must be in the format abc-abc
	re := regexp.MustCompile(`[a-z]+-[a-z]`)
	sku := re.FindAllString(fl.Field().String(), -1)

	if len(sku) == 1 {
		return true
	}

	return false
}
