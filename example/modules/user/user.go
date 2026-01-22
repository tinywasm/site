package user

type User struct {
	ID   int
	Name string
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
