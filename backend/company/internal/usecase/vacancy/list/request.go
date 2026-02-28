package list_vacancy

type Request struct {
	Order         Order
	Limit         int
	EncodedCursor string
	Requirements  *Requirements
}
