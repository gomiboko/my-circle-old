package forms

import "mime/multipart"

type CircleForm struct {
	CircleName     string                `form:"circleName" binding:"required,max=45,notonlywhitespace"`
	CircleIconFile *multipart.FileHeader `form:"circleIconFile"`
}
