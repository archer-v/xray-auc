package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var v *validator.Validate

type UserRecord struct {
	Email    string `envconfig:"optional" validate:"required" `
	Password string `envconfig:"optional"`
	Flow     string `envconfig:"optional"`
	Cipher   string `envconfig:"optional"`
	Proto    string `envconfig:"optional"`
	Tag      string `validate:"required" envconfig:"optional"`
}

func (u UserRecord) V() {

}
func (u UserRecord) Validate() (err error) {
	if v == nil {
		v = validator.New(validator.WithRequiredStructEnabled())
	}
	err = v.Struct(u)
	if err != nil {
		fmt.Printf("Some user options are incorrectly configured for user %v: ", u.Email)
		for _, e := range err.(validator.ValidationErrors) {
			if e.Tag() == "required" {
				fmt.Printf("	%v is required \n", strings.ToLower(e.StructField()))
			} else if e.Tag() == "oneof" {
				fmt.Printf("	%v should be one of: %v\n", strings.ToLower(e.StructField()), e.Param())
			} else {
				fmt.Printf("	%v: %v\n", strings.ToLower(e.StructField()), e)
			}
		}
	}
	return
}
