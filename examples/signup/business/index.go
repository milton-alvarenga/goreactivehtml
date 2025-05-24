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

func main (){

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
	ErrorMsg = "Check your inbox mail to confirm your account creation on the validation link sent on the message"
}