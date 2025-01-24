package internal

import "go.uber.org/fx"

// want "module.go files should not be directly in internal/ or internal/module/ directories"
var Module = fx.Module("internal",
	fx.Provide(
		NewInternalService,
	),
)
