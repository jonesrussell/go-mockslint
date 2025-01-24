package bad

import "go.uber.org/fx"

// want "module name \"wrong\" should match package name \"bad\""
var Module = fx.Module("wrong",
	fx.Provide(
		NewService,
	),
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
