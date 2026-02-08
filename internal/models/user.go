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

func (u *User) IsVerified() bool {
	if val, ok := u.Data["verified"].(bool); ok {
		return val
	}
	// SQLite might return 0/1 for bool
	if val, ok := u.Data["verified"].(int64); ok {
		return val == 1
	}
	return false
}
