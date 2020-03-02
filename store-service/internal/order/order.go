package order

import (
	"log"
	"store-service/internal/product"
)

type OrderService struct {
	ProductRepository product.ProductRepository
	OrderRepository   OrderInterface
}

type ProductRepository interface {
	GetProductByID(id int) product.Product
}

func (orderService OrderService) CreateOrder(submitedOrder SubmitedOrder) Order {
	totalAmount := orderService.GetTotalAmount(submitedOrder)

	orderID, err := orderService.OrderRepository.CreateOrder(totalAmount)
	if err != nil {
		log.Printf("OrderRepository.CreateOrder internal error %s", err.Error())
		return Order{}
	}

	shippingInfo := ShippingInfo{
		ShippingMethod:       submitedOrder.ShippingMethod,
		ShippingAddress:      submitedOrder.ShippingAddress,
		ShippingSubDistrict:  submitedOrder.ShippingSubDistrict,
		ShippingDistrict:     submitedOrder.ShippingDistrict,
		ShippingProvince:     submitedOrder.ShippingProvince,
		ShippingZipCode:      submitedOrder.ShippingZipCode,
		RecipientName:        submitedOrder.RecipientName,
		RecipientPhoneNumber: submitedOrder.RecipientPhoneNumber,
	}
	err = orderService.OrderRepository.CreatedShipping(orderID, shippingInfo)
	if err != nil {
		log.Printf("OrderRepository.CreatedShipping internal error %s", err.Error())
		return Order{}
	}

	for _, selectedProduct := range submitedOrder.Cart {
		err = orderService.OrderRepository.CreatedOrderProduct(orderID, selectedProduct.ProductID, selectedProduct.Quantity, selectedProduct.ProductPrice)
		if err != nil {
			log.Printf("OrderRepository.CreatedOrderProduct internal error %s", err.Error())
			return Order{}
		}
	}
	return Order{}
}

func (orderService OrderService) GetTotalProductPrice(submitedOrder SubmitedOrder) float64 {
	totalProductPrice := 0.00
	for _, cartItem := range submitedOrder.Cart {
		product, _ := orderService.ProductRepository.GetProductByID(cartItem.ProductID)
		totalProductPrice += product.Price * float64(cartItem.Quantity)
	}
	return totalProductPrice
}

func (orderService OrderService) GetTotalAmount(order SubmitedOrder) float64 {
	return orderService.GetTotalProductPrice(order) + order.GetShippingFee()
}
