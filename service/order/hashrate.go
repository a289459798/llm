package order

type HashRateOrder struct {
}

func NewHashRate(data OrderData) *HashRateOrder {
	return &HashRateOrder{}
}

func (hashRateOrder *HashRateOrder) Create(order CreateRequest) (response CreateResponse, err error) {
	return
}

func (hashRateOrder *HashRateOrder) Pay(orderData PayRequest) error {
	return nil
}
