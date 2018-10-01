package main

import (
	"github.com/cahlchang/lambda-go-processor/libs"
)

func main() {
	libs.Callee(libs.LambdaEvent{Response: ""})
}
