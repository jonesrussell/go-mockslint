package auth

import "go.uber.org/fx"

// GOOD: Module name matches package name and is in correct location
var Module = fx.Module("auth",
	fx.Provide(
		NewAuthenticator,
	),
)
