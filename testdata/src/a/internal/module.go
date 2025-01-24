package internal

import "go.uber.org/fx"

var Module = fx.Module("internal", // want "module.go files should not be directly in internal/ or internal/module/ directories"
	fx.Provide(
		NewInternalService,
	),
)
