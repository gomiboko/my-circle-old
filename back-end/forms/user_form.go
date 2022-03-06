package forms

type UserForm struct {
	Username string `json:"username" binding:"required,max=45,notonlywhitespace"`
	Email    string `json:"email"    binding:"required,max=254,email"`
	Password string `json:"password" binding:"required,min=8,max=128,password"`
}
