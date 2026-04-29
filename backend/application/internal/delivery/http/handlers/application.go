package handlers

type Application struct {
	create any
}

func NewApplication() *Application {
	return &Application{}
}
