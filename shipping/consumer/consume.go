package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"shipping/core"
	"shipping/logger"
)

type Port struct {
	Sub     core.Subscriber
	Service core.ShippingService
	Log     *logger.Logger
}

func (p *Port) Run(ctx context.Context) {

	rec := make(chan *core.Message)

	go p.Sub.Subscribe(ctx, "shipping", rec)

	for r := range rec {

		switch r.Topic {

		case "shipping":
			fmt.Println("SHIPPING", string(r.Msg))

			var ship *core.Shipping

			if err := json.Unmarshal(r.Msg, &ship); err != nil {
				p.Log.Println("ERROR:consumer.RUN.shipping", err)
			}

			fmt.Println(ship)

			if _, err := p.Service.AddOrderShipping(ctx, ship); err != nil {
				p.Log.Println("ERROR:consumer.RUN.shipping.AddOrderShipping", err)
			}
		}
	}
}
