package location

type (
	Order struct {
		OrderUUID string
	}
)

func NewOrder(orderUUID string) Order {
	return Order{OrderUUID: orderUUID}
}
