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
	el, _ := dom.Get("user-list")
	el.SetHTML("<p>Lista cargada via WASM (isomorfismo total)</p>")
}

func (u *User) OnUnmount() {}
