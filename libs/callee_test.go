package libs

import (
	"testing"
	"reflect"
)

func TestSanityReadConfig(t *testing.T) () {
	lstExpectationCaller := &ListCaller{
		Caller: []Caller{
			Caller{
				Name: "echo",
				Args: []string{"-e", "one\ntwo\nthree"},
				Pipe: []Pipe{
					Pipe{
						Name: "grep",
						Args: []string{"two"},
					},
				},
			},
		},
	}

	lstDummyCaller := &ListCaller{
		Caller: []Caller{
			Caller{
				Name: "echo",
				Args: []string{"-e", "four"},
				Pipe: []Pipe{
					Pipe{
						Name: "grep",
						Args: []string{"two"},
					},
				},
			},
		},
	}
	
	lstRes, err := ReadConfig("../testconf/accurate_config.yml")

	if err != nil {
		 t.Errorf("err found\n%v%v", lstRes, err)
	 }
	
	if !reflect.DeepEqual(lstExpectationCaller, lstRes) {
		t.Errorf("Unexpected value\nExpectation%v, Response %v", lstExpectationCaller, lstRes)
	}
	
	if reflect.DeepEqual(lstDummyCaller, lstRes) {
		t.Errorf("Unexpected value\nExpectation%v, Response %v", lstExpectationCaller, lstRes)
	}
}


func TestInvalidReadConfig(t *testing.T) () {
	res, err := ReadConfig("../testconf/invalid_config.yml")
	
	if err == nil {
		t.Errorf("no error\n%v%v", res, err)
	 }
}

func TestSanityProcessing(t *testing.T) () {
	msgSanity := "two\n"
	lstSanityCaller := &ListCaller{
		Caller: []Caller{
			Caller{
				Name: "echo",
				Args: []string{"-e", "one\ntwo\nthree"},
				Pipe: []Pipe{
					Pipe{
						Name: "grep",
						Args: []string{"two"},
					},
				},
			},
		},
	}
	msgRes := Processing(lstSanityCaller)
	if msgSanity !=  msgRes {
		t.Errorf("Unexpected value\nExpectation%v, Response %v", msgRes, msgSanity)
	}
}

func TestFailProcessing(t *testing.T) () {
	msgFail := `"/bin/expr"
"expr: division by zero\n"`
	lstSanityCaller := &ListCaller{
		Caller: []Caller{
			Caller{
				Name: "expr",
				Args: []string{"1", "/", "0"},
				Pipe: []Pipe{},
			},
		},
	}

	msgRes := Processing(lstSanityCaller)
	if msgFail !=  msgRes {
		t.Errorf("Unexpected value\nExpectation'%s', Response '%s'", msgFail, msgRes)
	}		
}
