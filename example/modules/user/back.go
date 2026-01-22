//go:build !wasm

package user

func (u *User) Read(data ...any) any {
	// Mock users
	return []*User{
		{ID: 1, Name: "Cesar"},
		{ID: 2, Name: "Admin"},
	}
}
