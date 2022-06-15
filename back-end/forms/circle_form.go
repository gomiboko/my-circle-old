package forms

type CircleForm struct {
	CircleName string `json:"circleName" binding:"required,max=45,notonlywhitespace"`
}
