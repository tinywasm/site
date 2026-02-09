package user

import "github.com/tinywasm/fmt"

type User struct {
	ID    int
	Name  string
	Param string
}

func (u *User) HandlerName() string {
	return "users"
}

func (u *User) DisplayName() string {
	return "Usuarios"
}

func (u *User) AllowedRoles(action byte) []byte {
	return []byte{'*'}
}

func (u *User) ValidateData(action byte, data ...any) error {
	return nil
}

func Add() []any {
	return []any{&User{}}
}

// Implement Parameterized interface
func (u *User) SetParams(params []string) {
	if len(params) > 0 {
		fmt.Println("User SetParams:", params)
		u.Param = params[0]
	} else {
		u.Param = ""
	}
}

// Implement ModuleLifecycle interface
func (u *User) BeforeNavigateAway() bool {
	fmt.Println("User BeforeNavigateAway")
	return true
}

func (u *User) AfterNavigateTo() {
	fmt.Println("User AfterNavigateTo")
	if u.Param != "" {
		fmt.Println("User navigated with param:", u.Param)
	}
}
