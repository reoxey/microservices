// Author: reoxey
// Date: 28-07-2021 09:15

package consumer

import "payment/core"

type Port struct {
	Sub     core.Subscriber
	Service core.PayService
}
