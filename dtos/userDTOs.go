package dtos

type CreateAccount struct {
	Username string
	Password string
	Email    string
}

type LoginAccount struct {
	Username string
	Password string
}
