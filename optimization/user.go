package main

// easyjson:json
type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func (u *User) Reset() {
	u.Name = ""
	u.Email = ""
	u.Browsers = nil
}
