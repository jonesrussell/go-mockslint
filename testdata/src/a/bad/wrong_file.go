package bad

import "go.uber.org/fx"

func init() { // want "fx.Module can only be used in module.go files"
	fx.Module("bad",
		fx.Provide(
			NewWrongService,
		),
	)
}

type WrongService struct{}

func NewWrongService() *WrongService {
	return &WrongService{}
}
