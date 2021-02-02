package profile

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
	IsAdmin  bool	`json:"is_admin"`
	JoinedAt string `json:"joined_at"`
}

type Users []User

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
