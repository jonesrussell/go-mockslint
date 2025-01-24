package auth

import (
	auth "a/internal/auth"

	"go.uber.org/fx"
)

// want "module.go files should not be directly in internal/ or internal/module/ directories"
// want "module name \"auth\" should match directory name \"module\""
var Module = fx.Module("auth",
	fx.Provide(
		auth.NewAuthenticator,
	),
)
