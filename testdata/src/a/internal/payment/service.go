package payment

import "go.uber.org/fx"

func init() { // want "fx.Module can only be used in module.go files"
	fx.Module("payment",
		fx.Provide(
			NewPaymentService,
		),
	)
}
