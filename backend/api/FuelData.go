package api

type FuelsData interface {
	GetFuels(middleware Middleware) error
}
