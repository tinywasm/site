package modules

import (
	"github.com/tinywasm/site/example/modules/contact"
	"github.com/tinywasm/site/example/modules/user"
)

func Init() []any {
	var all []any
	all = append(all, contact.Add()...)
	all = append(all, user.Add()...)
	return all
}
