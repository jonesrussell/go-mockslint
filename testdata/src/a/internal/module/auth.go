package auth

import (
	auth "a/internal/auth"

	"go.uber.org/fx"
)

// want "module.go files should not be directly in internal/ or internal/module/ directories"
var Module = fx.Module("auth", // want "module name \"auth\" should match directory name \"module\""
	fx.Provide(
		auth.NewAuthenticator,
	),
)
