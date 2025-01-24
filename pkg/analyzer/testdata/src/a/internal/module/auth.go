package auth

import "go.uber.org/fx"

// want "module.go files should not be directly in internal/ or internal/module/ directories"
var Module = fx.Module("auth",
	fx.Provide(
		NewAuthenticator,
	),
)
