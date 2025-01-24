package user

import "go.uber.org/fx"

var Module = fx.Module("auth", // want "module name \"auth\" should match package name \"user\""
	fx.Provide(
		NewUser,
	),
)
