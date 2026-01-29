package modules

import (
	"github.com/tinywasm/site"
	"github.com/tinywasm/site/example/modules/contact"
	"github.com/tinywasm/site/example/modules/user"
)

func Init() []any {
	// Configure Security (shared between Front and Back)
	site.SetUserRoles(func(data ...any) []byte {
		return []byte{'*'}
	})

	return append(
		contact.Add(),
		user.Add()...,
	)
}
