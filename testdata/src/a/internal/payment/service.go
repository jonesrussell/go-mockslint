package payment

import "go.uber.org/fx"

// want "fx.Module can only be used in module.go files"
func init() {
	fx.Module("payment",
		fx.Provide(
			NewPaymentService,
		),
	)
}
