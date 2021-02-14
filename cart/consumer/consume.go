package consumer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cart/core"
	"cart/logger"
)

type Port struct {
	Sub     core.Subscriber
	Service core.CartService
	Log     *logger.Logger
}

func (p *Port) Run(ctx context.Context) {

	rec := make(chan *core.Message)

	go p.Sub.Subscribe(ctx, "product_price", rec)

	for r := range rec {

		switch r.Topic {

		case "product_price":
			fmt.Println("PRODUCT_PRICE", string(r.Msg))

			arr := strings.Split(string(r.Msg), "|")

			id, err := strconv.Atoi(arr[0])
			if err != nil {
				p.Log.Println("ERROR:consumer.Run.product", err)
			}
			price, err := strconv.ParseFloat(arr[1], 64)
			if err != nil {
				p.Log.Println("ERROR:consumer.Run.product", err)
			}

			err = p.Service.UpdateItems(ctx, &core.Item{
				Id:    id,
				Price: price,
			})
			if err != nil {
				p.Log.Println("ERROR:consumer.Run.product", err)
			}
		}
	}
}
