package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"order/core"
	"order/logger"
)

type Port struct {
	Sub     core.Subscriber
	Service core.OrderService
	Log     *logger.Logger
}

func (p *Port) Run(ctx context.Context) {

	rec := make(chan *core.Message)

	go p.Sub.Subscribe(ctx, "order", rec)

	for r := range rec {

		switch r.Topic {

		case "order":
			fmt.Println("ORDER", string(r.Msg))

			cartOrder := &core.CartOrder{}

			if err := json.Unmarshal(r.Msg, &cartOrder); err != nil {
				p.Log.Println("ERROR:consumer.RUN.order", err)
				continue
			}

			payHandle := 0
			payStatus := 1
			orderStatus := core.OrderPaid
			if cartOrder.Checkout.PaymentMethod == core.PaymentCOD {
				payHandle = 1
				payStatus = 0
				orderStatus = core.OrderCreated
			}

			order := &core.Order{
				CartId:  cartOrder.Id,
				Buyer:   cartOrder.Buyer,
				Items:   cartOrder.Items,
				Payment: cartOrder.Payment,
				Status:  orderStatus,
			}

			id, err := p.Service.PlaceOrder(ctx, order)
			if err != nil {
				p.Log.Println("ERROR:consumer.RUN.order.PlaceOrder", err)
				continue
			}

			go func() {
				shipping := &core.Shipping{
					OrderId:       id,
					Address: 	   &core.Address{Id: cartOrder.Checkout.AddressId},
					Status:        core.Ordered,
					PaymentHandle: payHandle,
					PaymentStatus: payStatus,
				}

				if err = p.Service.OrderShipping(ctx, shipping); err != nil {
					p.Log.Println("ERROR:consumer.RUN.order.OrderShipping", err)
				}
			}()

			go func() {

				var items []*core.Item
				for _, item := range order.Items {
					items = append(items, &core.Item{Id: item.Id, Qty: item.Qty})
				}

				if err = p.Service.UpdateItemStocks(ctx, items); err != nil {
					p.Log.Println("ERROR:consumer.RUN.order.UpdateItemStocks", err)
				}
			}()
		}
	}
}
