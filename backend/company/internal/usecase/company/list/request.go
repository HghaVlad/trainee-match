package list

type Request struct {
	Order  Order
	Limit  int
	Cursor string
}
