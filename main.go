package main

import (
	"github.com/cahlchang/lambda-go-processor/libs"
	//"github.com/aws/aws-lambda-go/lambda"
)


func main() {
	//lambda.Start(libs.Callee)
	libs.Callee(libs.LambdaEvent{Response: ""})
}
