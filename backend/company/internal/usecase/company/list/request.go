package list_companies

type Request struct {
	Order  Order
	Limit  int
	Cursor string
}
