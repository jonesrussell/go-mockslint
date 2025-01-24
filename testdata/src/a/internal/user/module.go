package user

import "go.uber.org/fx"

// want "module name \"auth\" should match package name \"user\""
var Module = fx.Module("auth",
	fx.Provide(
		NewUser,
	),
)
