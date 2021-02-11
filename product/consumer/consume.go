package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"product/core"
	"product/logger"
)

type Port struct {
	Sub     core.Subscriber
	Service core.ProductService
	Log     *logger.Logger
}

func (p *Port) Run(ctx context.Context) {

	rec := make(chan *core.Message)

	go p.Sub.Subscribe(ctx, "product_qty", rec)

	for r := range rec {

		switch r.Topic {

		case "product_qty":
			fmt.Println("PRODUCT_QTY", string(r.Msg))

			var items core.Items

			if err := json.Unmarshal(r.Msg, &items); err != nil {
				p.Log.Println("ERROR:consumer.RUN.product_qty", err)
			}

			if err := p.Service.UpdateProductStocks(ctx, items); err != nil {
				p.Log.Println("ERROR:consumer.RUN.product_qty.UpdateProductStocks", err)
			}

		}
	}
}
