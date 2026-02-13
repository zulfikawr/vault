package models

type User struct {
	*Record
}

func AsUser(r *Record) *User {
	if r == nil {
		return nil
	}
	return &User{Record: r}
}

func (u *User) Username() string {
	return u.GetString("username")
}

func (u *User) Email() string {
	return u.GetString("email")
}
