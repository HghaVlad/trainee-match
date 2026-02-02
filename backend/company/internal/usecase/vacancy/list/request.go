package list_vacancy

type Request struct {
	Order  Order
	Limit  int
	Cursor string
}
