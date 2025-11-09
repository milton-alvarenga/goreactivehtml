package main

import (
	"github.com/milton-alvarenga/goreactivehtml/examples/signup/business"
)

func main() {
	business.SubmitSignup("argEmail", "argPassword", "argConfirmPassword")
}
