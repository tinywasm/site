//go:build wasm

package user

import (
	"github.com/tinywasm/dom"
)

func (u *User) RenderHTML() string {
	return `<!-- module -->
<section id="users">
    <h1>Usuarios</h1>
    <div id="user-list">Cargando...</div>
</section>`
}

func (u *User) OnMount() {
	dom.Render("user-list", dom.P("Lista cargada via WASM (isomorfismo total)"))
}

func (u *User) OnUnmount() {}

func (u *User) GetID() string             { return "user-module" }
func (u *User) SetID(id string)           {}
func (u *User) Children() []dom.Component { return nil }
