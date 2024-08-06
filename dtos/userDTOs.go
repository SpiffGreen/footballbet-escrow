package dtos

type CreateAccount struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
}

type LoginAccount struct {
	Username string
	Password string
}
