package modules

import (
	"github.com/tinywasm/site/example/modules/contact"
	"github.com/tinywasm/site/example/modules/user"
)

func Init() []any {
	return append(
		contact.Add(),
		user.Add()...,
	)
}
