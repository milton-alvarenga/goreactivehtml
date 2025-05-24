package business

import(
	"database/sql"

	_ "github.com/lib/pq"

	"projectName/util/validation"
)

var Email string
var Password string
var ConfirmPassword string
var ErrorMsg string
var SuccessMsg string

func main (){

}

type OutputMsg struct {

}

func SubmitSignup(Email string, Password string, ConfirmPassword string){
	//Checa
	//Processa
	
	if ! validation.Email(Email){
		ErrorMsg = "Invalid email"
		return
	}

	Email = ""
	Password = ""
	ConfirmPassword == ""
	ErrorMsg = ""
	SuccessMsg = "Check your inbox mail to confirm your account creation on the validation link sent on the message"
}

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