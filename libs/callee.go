package libs

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"bytes"
	"log"
	yaml "gopkg.in/yaml.v2"
)

// ListCaller is a structure containing a list of call information.
type ListCaller struct {
	Caller []Caller `yaml:"caller"`
}

// Caller is a structure containing information on the process to be executed.
type Caller struct {
	Name             string           `yaml:"name"`
	Args               []string         `yaml:"args"`
	Pipe               []Pipe              `yaml:"pipe"`
}

// Pipe is a structure containing pipeline information.
type Pipe struct {	
	Name       string           `yaml:"name"`
	Args        []string         `yaml:"args"`
}

// LambdaEvent is a structure required by AWS LAMBDA's SDK.
type LambdaEvent struct {
	Response string `json:"dummy event"`
}

// LambdaResponse is a structure required by AWS LAMBDA's SDK.
type LambdaResponse struct {
	Message string `json:"Response:"`
}

// ReadConfig function reads the configuration defined in YMAL format.
func ReadConfig(pathYaml string) (*ListCaller, error) {
	buf, errRead := ioutil.ReadFile(pathYaml)

	if errRead != nil {
		return  nil, errRead
	}

	var l ListCaller
	errMarchal := yaml.Unmarshal(buf, &l)
	if errMarchal != nil {
		return nil, errMarchal
	}
	return &l, nil
}

// Processing function executes the process as defined in the structure.
func Processing(l *ListCaller)(string) {
	msgRet := ""

	for _, caller := range l.Caller {
		cmdMain := exec.Command(caller.Name, caller.Args[0:]... )
		cmdPre := cmdMain
		
		lstCmd := make([]*exec.Cmd, 0)
		lstCmd = append(lstCmd, cmdPre)

		// Connecting the output of the execution process of pipe connection with the previous process.
		for _, pipe := range caller.Pipe {
			cmdPipe := exec.Command(pipe.Name, pipe.Args[0:]... )
			cmdPipe.Stdin, _= cmdPre.StdoutPipe()
			cmdPre = cmdPipe
			
			lstCmd = append(lstCmd, cmdPipe)
		}

		// After pipeline connection, secure output to buffer.
		var outbufMain  bytes.Buffer
		var outbufErr  bytes.Buffer
		
		lstCmd[len(lstCmd)-1].Stdout = &outbufMain
		lstCmd[len(lstCmd)-1].Stderr = &outbufErr
		
		for i := range lstCmd {
			lstCmd[len(lstCmd) - i - 1].Start()
		}
		
		for _, cmd := range lstCmd {
			e := cmd.Wait()
			if e != nil {
				msgRet += tee("%q\n%q", cmd.Path, outbufErr.String())
			} 
		}
		fmt.Printf("%q command end\n", caller.Name)
		msgRet += tee(outbufMain.String())
	}
	
	fmt.Printf("all command end\n")
	return msgRet
}

// Callee function is called from the AWS SDK.
func Callee(event LambdaEvent) (LambdaResponse, error) {
	event.Response = ""

	l, err := ReadConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
		return  LambdaResponse{Message: "panic config"}, err
	}
	return LambdaResponse{Message: fmt.Sprintf("%s", Processing(l))}, nil
}

func tee(msgAddFmt string, msgArgs ...interface{}) (string) {
	var msg string
	if 0 != len(msgArgs) {
		msg = fmt.Sprintf(msgAddFmt, msgArgs...)
	} else {
		
		msg = fmt.Sprint(msgAddFmt)
	}

	fmt.Print(msg)
	return msg
}
