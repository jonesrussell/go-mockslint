package a

import "go.uber.org/fx"

// This is a valid module
var Module = fx.Module("a",
	fx.Provide(
		NewService,
	),
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
