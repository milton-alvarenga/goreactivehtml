package business

import (
	"fmt"
)

var Email string
var Password string
var ConfirmPassword string
var ErrorMsg string
var SuccessMsg string

func UpdateDescription(email string, login string) string {
	return Email
}

func SubmitSignup(jsEmail string, jsPassword string, jsConfirmPassword string) {
	//Checa
	//Processa

	//if !validation.Email(Email) {
	if false {
		ErrorMsg = "Invalid email"
		return
	}
	fmt.Println("Before")
	fmt.Println("Global email", jsEmail)
	fmt.Println("Global password", jsPassword)
	fmt.Println("Global confirmpassword", jsConfirmPassword)
	fmt.Println("Global successMsg", SuccessMsg)
	fmt.Println("Local email", Email)
	fmt.Println("Local password", Password)
	fmt.Println("Local confirmpassword", ConfirmPassword)
	fmt.Println("Local successMsg", SuccessMsg)

	Email = ""
	Password = ""
	ConfirmPassword = ""
	ErrorMsg = ""
	SuccessMsg = "Check your inbox mail to confirm your account creation on the validation link sent on the message"

	fmt.Println("After")
	fmt.Println("Global email", jsEmail)
	fmt.Println("Global password", jsPassword)
	fmt.Println("Global confirmpassword", jsConfirmPassword)
	fmt.Println("Global successMsg", SuccessMsg)
	fmt.Println("Local email", Email)
	fmt.Println("Local password", Password)
	fmt.Println("Local confirmpassword", ConfirmPassword)
	fmt.Println("Local successMsg", SuccessMsg)
}

/*
transpilacao
setEmail("")
setPassword("")
setConfirmPassword("")
setErrorMsg("")
setSuccessMsg("")


set.Update(&Email,"")

updateJs(
	//verifca o dirty de cada variável para subir
	//as funcoes set quem alteram realmente o escopo global e sujaram o dirty
)
*/
